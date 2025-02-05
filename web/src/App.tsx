import { Routes, Route, useLocation } from "react-router";
import { lazy, Suspense } from "react";
import Loading from "@/components/Loading";
import Layout from "@/components/Layout";

// Lazy load pages
const Home = lazy(() => import("@/pages/Home"));
const Login = lazy(() => import("@/pages/Login"));
const Register = lazy(() => import("@/pages/Register"));
const ForgotPassword = lazy(() => import("@/pages/ForgotPassword"));
const ResetPassword = lazy(() => import("@/pages/ResetPassword"));

const Dashboard = lazy(() => import("@/pages/protected/Dashboard"));
const NotFound = lazy(() => import("@/pages/NotFound"));

export default function App() {
  const location = useLocation();

  return (
    <Suspense fallback={<Loading />}>
      <Routes location={location} key={location.pathname}>
        <Route index element={<Home />} />
        <Route path="login" element={<Login />} />
        <Route path="register" element={<Register />} />
        <Route path="forgot-password" element={<ForgotPassword />} />
        <Route path="reset-password" element={<ResetPassword />} />
        <Route element={<Layout />}>
          <Route path="dashboard" element={<Dashboard />} />
        </Route>
        <Route path="layout" element={<Layout />} />
        <Route path="*" element={<NotFound />} />
      </Routes>
    </Suspense>
  );
}
