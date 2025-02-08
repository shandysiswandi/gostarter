<script lang="ts">
  import { Router, Route } from "svelte-routing";
  import Home from "./pages/Home.svelte";
  import Login from "./pages/Login.svelte";
  import NotFound from "./pages/NotFound.svelte";
  import Register from "./pages/Register.svelte";
  import ResetPassword from "./pages/ResetPassword.svelte";
  import ForgotPassword from "./pages/ForgotPassword.svelte";
  import Dashboard from "./pages/Dashboard.svelte";
</script>

<Router>
  <Route path="/" component={Home} />
  <Route path="/login" component={Login} />
  <Route path="/register" component={Register} />
  <Route path="/forgot-password" component={ForgotPassword} />
  <Route path="/reset-password" component={ResetPassword} />
  <Route path="/dashboard" component={Dashboard} />
  <Route path="*" component={NotFound} />
</Router>
