import { lazy, Suspense } from "react";
import { Navigate, Route, Routes } from "react-router-dom";
import { AppLayout } from "./components/layout/AppLayout";
import { FullPageLoader } from "./components/ui/FullPageLoader";
import { AuthProvider } from "./context/AuthContext";
import { LoginPage } from "./pages/LoginPage";
import { NotFoundPage } from "./pages/NotFoundPage";
import { RegisterPage } from "./pages/RegisterPage";
import { ProtectedRoute } from "./routes/ProtectedRoute";

const DashboardPage = lazy(() =>
  import("./pages/DashboardPage").then((module) => ({ default: module.DashboardPage })),
);
const AddServicePage = lazy(() =>
  import("./pages/AddServicePage").then((module) => ({ default: module.AddServicePage })),
);
const EditServicePage = lazy(() =>
  import("./pages/EditServicePage").then((module) => ({ default: module.EditServicePage })),
);
const ServiceDetailPage = lazy(() =>
  import("./pages/ServiceDetailPage").then((module) => ({ default: module.ServiceDetailPage })),
);

function App() {
  return (
    <AuthProvider>
      <Suspense fallback={<FullPageLoader label="Loading page" />}>
        <Routes>
          <Route path="/login" element={<LoginPage />} />
          <Route path="/register" element={<RegisterPage />} />

          <Route element={<ProtectedRoute />}>
            <Route element={<AppLayout />}>
              <Route path="/dashboard" element={<DashboardPage />} />
              <Route path="/services/new" element={<AddServicePage />} />
              <Route path="/services/:id" element={<ServiceDetailPage />} />
              <Route path="/services/:id/edit" element={<EditServicePage />} />
            </Route>
          </Route>

          <Route path="/" element={<Navigate to="/dashboard" replace />} />
          <Route path="*" element={<NotFoundPage />} />
        </Routes>
      </Suspense>
    </AuthProvider>
  );
}

export default App;
