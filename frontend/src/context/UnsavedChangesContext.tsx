import {
  createContext,
  useCallback,
  useContext,
  useRef,
  type ReactNode,
} from "react";

type NavigationGuard = () => boolean;

interface UnsavedChangesContextValue {
  registerNavigationGuard: (guard: NavigationGuard) => () => void;
  confirmNavigation: () => boolean;
}

const UnsavedChangesContext = createContext<UnsavedChangesContextValue | undefined>(undefined);

export function UnsavedChangesProvider({ children }: { children: ReactNode }) {
  const activeGuardRef = useRef<NavigationGuard | null>(null);

  const registerNavigationGuard = useCallback((guard: NavigationGuard) => {
    activeGuardRef.current = guard;

    return () => {
      if (activeGuardRef.current === guard) {
        activeGuardRef.current = null;
      }
    };
  }, []);

  const confirmNavigation = useCallback(
    () => activeGuardRef.current?.() ?? true,
    [],
  );

  return (
    <UnsavedChangesContext.Provider value={{ registerNavigationGuard, confirmNavigation }}>
      {children}
    </UnsavedChangesContext.Provider>
  );
}

// oxlint-disable-next-line react/only-export-components -- colocating the hook keeps the guard API cohesive.
export function useUnsavedChanges() {
  const context = useContext(UnsavedChangesContext);

  if (!context) {
    throw new Error("useUnsavedChanges must be used inside UnsavedChangesProvider");
  }

  return context;
}
