import { Link } from "react-router";
import { useState } from "react";
import { useDocumentTitle } from "@/hooks";

const ForgotPassword: React.FC = () => {
  useDocumentTitle("Forgot Password");

  const [email, setEmail] = useState("");
  const [isLoading, setIsLoading] = useState(false);
  const [isSubmitted, setIsSubmitted] = useState(false);

  const handleSubmit = (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();

    setIsSubmitted(false);
    setIsLoading(true);
    setEmail("");

    setTimeout(function () {
      // Simulating API call
      setIsLoading(false);
      setIsSubmitted(true);
    }, 1000);

    // await api call to send reset password link
  };

  return (
    <div className="flex min-h-screen items-center justify-center bg-gray-100 p-4">
      <div className="w-full max-w-md rounded-xl bg-white p-8 shadow-lg">
        <h2 className="mb-2 text-center text-2xl font-bold text-gray-900">
          Forgot Password
        </h2>
        <h4 className="mb-6 text-center text-sm text-gray-600">
          We will send password reset link on your email
        </h4>

        {isSubmitted && (
          <div className="mb-4 rounded-lg bg-green-50 p-4 text-sm text-green-700">
            We have sent you an email with instructions to reset your password.
          </div>
        )}

        <form className="space-y-4" onSubmit={handleSubmit}>
          <div>
            <label
              htmlFor="email"
              className="mb-1 block text-sm font-medium text-gray-700"
            >
              Email
            </label>
            <input
              type="email"
              id="email"
              name="email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              className="w-full rounded-lg border border-gray-300 px-4 py-2 outline-none transition-all focus:border-indigo-500 focus:ring-2 focus:ring-indigo-500"
              placeholder="your@email.com"
            />
          </div>

          <button
            type="submit"
            disabled={isLoading} // Disable button during loading
            className={`w-full rounded-lg py-2.5 font-medium text-white transition-colors ${
              isLoading
                ? "cursor-not-allowed bg-indigo-400"
                : "bg-indigo-600 hover:bg-indigo-700"
            }`}
          >
            {isLoading ? "Loading..." : "Send Reset Instructions"}
          </button>
        </form>

        <div className="mt-6 text-center text-sm text-gray-600">
          Don't have an account?
          <Link
            to="/register"
            className="font-medium text-indigo-600 hover:text-indigo-500"
          >
            Sign up
          </Link>
        </div>
      </div>
    </div>
  );
};

export default ForgotPassword;
