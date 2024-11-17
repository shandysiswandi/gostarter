import grpc from "k6/net/grpc";
import { faker } from "https://esm.sh/@faker-js/faker";
import { check } from "k6";

const client = new grpc.Client();
client.load(["./proto"], "auth.proto");

export default function () {
  client.connect("localhost:50001", { plaintext: true });

  const pathAuth = "gostarter.api.auth.AuthService";

  const login = client.invoke(`${pathAuth}/Login`, {
    email: "admin@admin.com",
    password: "admin123",
  });
  check(login, {
    "login is success": (r) => r && r.status === grpc.StatusOK,
  });

  const email = faker.internet.email();
  const register = client.invoke(`${pathAuth}/Register`, {
    email: email,
    password: email,
  });
  check(register, {
    "register is success": (r) => r && r.status === grpc.StatusOK,
  });

  const refreshToken = client.invoke(`${pathAuth}/RefreshToken`, {
    refresh_token: login.message.refreshToken,
  });
  check(refreshToken, {
    "refresh token is success": (r) => r && r.status === grpc.StatusOK,
  });

  const forgotPassword = client.invoke(`${pathAuth}/ForgotPassword`, {
    email: email,
  });
  check(forgotPassword, {
    "forgot password is success": (r) => r && r.status === grpc.StatusOK,
  });

  const resetPassword = client.invoke(`${pathAuth}/ResetPassword`, {
    token: email,
    password: email,
  });
  check(resetPassword, {
    "reset password is success": (r) => r && r.status === grpc.StatusOK,
  });

  client.close();
}
