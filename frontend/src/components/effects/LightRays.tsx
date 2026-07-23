import { useEffect, useRef } from "react";
import { usePrefersReducedMotion } from "../../hooks/usePrefersReducedMotion";

type RayOrigin =
  | "top-left"
  | "top-center"
  | "top-right"
  | "left"
  | "right"
  | "bottom-left"
  | "bottom-center"
  | "bottom-right";

interface LightRaysProps {
  className?: string;
  color?: string;
  origin?: RayOrigin;
  speed?: number;
  spread?: number;
  length?: number;
  fadeDistance?: number;
  saturation?: number;
  mouseInfluence?: number;
  noise?: number;
  distortion?: number;
}

const vertexShader = `
attribute vec2 position;

void main() {
  gl_Position = vec4(position, 0.0, 1.0);
}`;

// Adapted from Background Animations / 1 in the supplied Awwwards Pack.
const fragmentShader = `
precision highp float;

uniform float iTime;
uniform vec2 iResolution;
uniform vec2 rayPos;
uniform vec2 rayDir;
uniform vec3 raysColor;
uniform float raysSpeed;
uniform float lightSpread;
uniform float rayLength;
uniform float fadeDistance;
uniform float saturation;
uniform vec2 mousePos;
uniform float mouseInfluence;
uniform float noiseAmount;
uniform float distortion;

float randomNoise(vec2 point) {
  return fract(sin(dot(point.xy, vec2(12.9898, 78.233))) * 43758.5453123);
}

float rayStrength(
  vec2 source,
  vec2 referenceDirection,
  vec2 coordinate,
  float seedA,
  float seedB,
  float speed
) {
  vec2 sourceToCoordinate = coordinate - source;
  vec2 direction = normalize(sourceToCoordinate);
  float angle = dot(direction, referenceDirection);
  float distortedAngle = angle
    + distortion * sin(iTime * 2.0 + length(sourceToCoordinate) * 0.01) * 0.2;
  float spreadFactor = pow(max(distortedAngle, 0.0), 1.0 / max(lightSpread, 0.001));
  float distanceFromSource = length(sourceToCoordinate);
  float maximumDistance = iResolution.x * rayLength;
  float lengthFalloff = clamp(
    (maximumDistance - distanceFromSource) / maximumDistance,
    0.0,
    1.0
  );
  float fadeFalloff = clamp(
    (iResolution.x * fadeDistance - distanceFromSource) / (iResolution.x * fadeDistance),
    0.5,
    1.0
  );
  float baseStrength = clamp(
    (0.45 + 0.15 * sin(distortedAngle * seedA + iTime * speed))
      + (0.3 + 0.2 * cos(-distortedAngle * seedB + iTime * speed)),
    0.0,
    1.0
  );

  return baseStrength * lengthFalloff * fadeFalloff * spreadFactor;
}

void main() {
  vec2 coordinate = vec2(gl_FragCoord.x, iResolution.y - gl_FragCoord.y);
  vec2 finalRayDirection = rayDir;

  if (mouseInfluence > 0.0) {
    vec2 mouseScreenPosition = mousePos * iResolution.xy;
    vec2 mouseDirection = normalize(mouseScreenPosition - rayPos);
    finalRayDirection = normalize(mix(rayDir, mouseDirection, mouseInfluence));
  }

  vec4 firstRays = vec4(1.0) * rayStrength(
    rayPos,
    finalRayDirection,
    coordinate,
    36.2214,
    21.11349,
    1.5 * raysSpeed
  );
  vec4 secondRays = vec4(1.0) * rayStrength(
    rayPos,
    finalRayDirection,
    coordinate,
    22.3991,
    18.0234,
    1.1 * raysSpeed
  );
  vec4 color = firstRays * 0.5 + secondRays * 0.4;

  if (noiseAmount > 0.0) {
    float noiseValue = randomNoise(coordinate * 0.01 + iTime * 0.1);
    color.rgb *= 1.0 - noiseAmount + noiseAmount * noiseValue;
  }

  float brightness = 1.0 - coordinate.y / iResolution.y;
  color.r *= 0.1 + brightness * 0.8;
  color.g *= 0.3 + brightness * 0.6;
  color.b *= 0.5 + brightness * 0.5;

  if (saturation != 1.0) {
    float gray = dot(color.rgb, vec3(0.299, 0.587, 0.114));
    color.rgb = mix(vec3(gray), color.rgb, saturation);
  }

  color.rgb *= raysColor;
  gl_FragColor = color;
}`;

function hexToRgb(hex: string) {
  const normalized = hex.trim().replace(/^#/, "");
  const expanded = normalized.length === 3
    ? normalized.split("").map((character) => `${character}${character}`).join("")
    : normalized;
  const match = /^([a-f\d]{2})([a-f\d]{2})([a-f\d]{2})$/i.exec(expanded);

  if (!match) return [1, 1, 1] as const;
  return [
    Number.parseInt(match[1], 16) / 255,
    Number.parseInt(match[2], 16) / 255,
    Number.parseInt(match[3], 16) / 255,
  ] as const;
}

function getAnchorAndDirection(origin: RayOrigin, width: number, height: number) {
  const outside = 0.2;

  switch (origin) {
    case "top-left":
      return { anchor: [0, -outside * height] as const, direction: [0, 1] as const };
    case "top-right":
      return { anchor: [width, -outside * height] as const, direction: [0, 1] as const };
    case "left":
      return { anchor: [-outside * width, 0.5 * height] as const, direction: [1, 0] as const };
    case "right":
      return { anchor: [(1 + outside) * width, 0.5 * height] as const, direction: [-1, 0] as const };
    case "bottom-left":
      return { anchor: [0, (1 + outside) * height] as const, direction: [0, -1] as const };
    case "bottom-center":
      return { anchor: [0.5 * width, (1 + outside) * height] as const, direction: [0, -1] as const };
    case "bottom-right":
      return { anchor: [width, (1 + outside) * height] as const, direction: [0, -1] as const };
    default:
      return { anchor: [0.5 * width, -outside * height] as const, direction: [0, 1] as const };
  }
}

export function LightRays({
  className = "",
  color = "#67e8f9",
  origin = "top-center",
  speed = 0.55,
  spread = 0.65,
  length = 2.2,
  fadeDistance = 1.15,
  saturation = 1.1,
  mouseInfluence = 0.08,
  noise = 0,
  distortion = 0.035,
}: LightRaysProps) {
  const containerRef = useRef<HTMLDivElement>(null);
  const canvasRef = useRef<HTMLCanvasElement>(null);
  const prefersReducedMotion = usePrefersReducedMotion();

  useEffect(() => {
    const container = containerRef.current;
    const canvas = canvasRef.current;
    if (!container || !canvas || prefersReducedMotion) return;

    const context = canvas.getContext("webgl", {
      alpha: true,
      antialias: false,
      powerPreference: "low-power",
      premultipliedAlpha: false,
    });
    if (!context) return;

    const compileShader = (type: number, source: string) => {
      const shader = context.createShader(type);
      if (!shader) return null;
      context.shaderSource(shader, source);
      context.compileShader(shader);
      if (!context.getShaderParameter(shader, context.COMPILE_STATUS)) {
        context.deleteShader(shader);
        return null;
      }
      return shader;
    };

    const compiledVertexShader = compileShader(context.VERTEX_SHADER, vertexShader);
    const compiledFragmentShader = compileShader(context.FRAGMENT_SHADER, fragmentShader);
    if (!compiledVertexShader || !compiledFragmentShader) {
      if (compiledVertexShader) context.deleteShader(compiledVertexShader);
      if (compiledFragmentShader) context.deleteShader(compiledFragmentShader);
      return;
    }

    const program = context.createProgram();
    if (!program) {
      context.deleteShader(compiledVertexShader);
      context.deleteShader(compiledFragmentShader);
      return;
    }
    context.attachShader(program, compiledVertexShader);
    context.attachShader(program, compiledFragmentShader);
    context.linkProgram(program);
    context.deleteShader(compiledVertexShader);
    context.deleteShader(compiledFragmentShader);
    if (!context.getProgramParameter(program, context.LINK_STATUS)) {
      context.deleteProgram(program);
      return;
    }

    context.useProgram(program);
    context.clearColor(0, 0, 0, 0);
    context.disable(context.DEPTH_TEST);
    context.disable(context.CULL_FACE);
    context.enable(context.BLEND);
    context.blendFunc(context.SRC_ALPHA, context.ONE_MINUS_SRC_ALPHA);

    const positionLocation = context.getAttribLocation(program, "position");
    const quad = context.createBuffer();
    if (positionLocation < 0 || !quad) {
      if (quad) context.deleteBuffer(quad);
      context.deleteProgram(program);
      return;
    }
    context.bindBuffer(context.ARRAY_BUFFER, quad);
    context.bufferData(
      context.ARRAY_BUFFER,
      new Float32Array([-1, -1, 3, -1, -1, 3]),
      context.STATIC_DRAW,
    );
    context.enableVertexAttribArray(positionLocation);
    context.vertexAttribPointer(positionLocation, 2, context.FLOAT, false, 0, 0);

    const locations = {
      time: context.getUniformLocation(program, "iTime"),
      resolution: context.getUniformLocation(program, "iResolution"),
      rayPosition: context.getUniformLocation(program, "rayPos"),
      rayDirection: context.getUniformLocation(program, "rayDir"),
      color: context.getUniformLocation(program, "raysColor"),
      speed: context.getUniformLocation(program, "raysSpeed"),
      spread: context.getUniformLocation(program, "lightSpread"),
      length: context.getUniformLocation(program, "rayLength"),
      fadeDistance: context.getUniformLocation(program, "fadeDistance"),
      saturation: context.getUniformLocation(program, "saturation"),
      mousePosition: context.getUniformLocation(program, "mousePos"),
      mouseInfluence: context.getUniformLocation(program, "mouseInfluence"),
      noise: context.getUniformLocation(program, "noiseAmount"),
      distortion: context.getUniformLocation(program, "distortion"),
    };

    const [red, green, blue] = hexToRgb(color);
    context.uniform3f(locations.color, red, green, blue);
    context.uniform1f(locations.speed, speed);
    context.uniform1f(locations.spread, spread);
    context.uniform1f(locations.length, length);
    context.uniform1f(locations.fadeDistance, fadeDistance);
    context.uniform1f(locations.saturation, saturation);
    context.uniform1f(locations.mouseInfluence, mouseInfluence);
    context.uniform1f(locations.noise, noise);
    context.uniform1f(locations.distortion, distortion);

    let resolutionWidth = 1;
    let resolutionHeight = 1;
    let rayPositionX = 0;
    let rayPositionY = 0;
    let rayDirectionX = 0;
    let rayDirectionY = 1;
    let mouseX = 0.5;
    let mouseY = 0.5;
    let smoothMouseX = 0.5;
    let smoothMouseY = 0.5;
    let animationFrame = 0;
    let inViewport = true;
    let documentVisible = !document.hidden;

    const updatePlacement = () => {
      const bounds = container.getBoundingClientRect();
      if (bounds.width <= 0 || bounds.height <= 0) return;
      const devicePixelRatio = Math.min(window.devicePixelRatio || 1, 1.75);
      resolutionWidth = Math.max(1, Math.floor(bounds.width * devicePixelRatio));
      resolutionHeight = Math.max(1, Math.floor(bounds.height * devicePixelRatio));

      if (canvas.width !== resolutionWidth || canvas.height !== resolutionHeight) {
        canvas.width = resolutionWidth;
        canvas.height = resolutionHeight;
        context.viewport(0, 0, resolutionWidth, resolutionHeight);
      }

      const placement = getAnchorAndDirection(origin, resolutionWidth, resolutionHeight);
      [rayPositionX, rayPositionY] = placement.anchor;
      [rayDirectionX, rayDirectionY] = placement.direction;
    };

    const draw = (time: number) => {
      if (!inViewport || !documentVisible) {
        animationFrame = 0;
        return;
      }

      const smoothing = 0.92;
      smoothMouseX = smoothMouseX * smoothing + mouseX * (1 - smoothing);
      smoothMouseY = smoothMouseY * smoothing + mouseY * (1 - smoothing);

      context.uniform1f(locations.time, time * 0.001);
      context.uniform2f(locations.resolution, resolutionWidth, resolutionHeight);
      context.uniform2f(locations.rayPosition, rayPositionX, rayPositionY);
      context.uniform2f(locations.rayDirection, rayDirectionX, rayDirectionY);
      context.uniform2f(locations.mousePosition, smoothMouseX, smoothMouseY);
      context.drawArrays(context.TRIANGLES, 0, 3);
      animationFrame = window.requestAnimationFrame(draw);
    };

    const startDrawing = () => {
      if (!animationFrame && inViewport && documentVisible) {
        animationFrame = window.requestAnimationFrame(draw);
      }
    };

    const stopDrawing = () => {
      if (!animationFrame) return;
      window.cancelAnimationFrame(animationFrame);
      animationFrame = 0;
    };

    const interactionSurface = container.parentElement ?? container;
    const handlePointerMove = (event: PointerEvent) => {
      const bounds = container.getBoundingClientRect();
      if (bounds.width <= 0 || bounds.height <= 0) return;
      mouseX = Math.min(1, Math.max(0, (event.clientX - bounds.left) / bounds.width));
      mouseY = Math.min(1, Math.max(0, (event.clientY - bounds.top) / bounds.height));
    };

    const handleVisibilityChange = () => {
      documentVisible = !document.hidden;
      if (documentVisible) startDrawing();
      else stopDrawing();
    };

    const resizeObserver = new ResizeObserver(updatePlacement);
    const intersectionObserver = new IntersectionObserver(([entry]) => {
      inViewport = entry.isIntersecting;
      if (inViewport) startDrawing();
      else stopDrawing();
    });

    updatePlacement();
    resizeObserver.observe(container);
    intersectionObserver.observe(container);
    interactionSurface.addEventListener("pointermove", handlePointerMove, { passive: true });
    document.addEventListener("visibilitychange", handleVisibilityChange);
    startDrawing();

    return () => {
      stopDrawing();
      resizeObserver.disconnect();
      intersectionObserver.disconnect();
      interactionSurface.removeEventListener("pointermove", handlePointerMove);
      document.removeEventListener("visibilitychange", handleVisibilityChange);
      context.deleteBuffer(quad);
      context.deleteProgram(program);
      const loseContext = context.getExtension("WEBGL_lose_context");
      loseContext?.loseContext();
    };
  }, [
    color,
    distortion,
    fadeDistance,
    length,
    mouseInfluence,
    noise,
    origin,
    prefersReducedMotion,
    saturation,
    speed,
    spread,
  ]);

  return (
    <div ref={containerRef} className={`light-rays-layer ${className}`} aria-hidden="true">
      <canvas ref={canvasRef} />
    </div>
  );
}
