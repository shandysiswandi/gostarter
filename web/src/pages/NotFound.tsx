import { Link } from "react-router";
import { useDocumentTitle } from "@/hooks";

const NotFound: React.FC = () => {
  useDocumentTitle("Not Found");

  return (
    <div className="flex h-screen items-center bg-gray-100 text-gray-800">
      <div className="container mx-auto flex flex-col items-center justify-center">
        <div className="max-w-md text-center">
          <h2 className="mb-8 text-9xl font-extrabold text-gray-400">
            <span className="sr-only">Error</span>404
          </h2>
          <p className="text-2xl font-semibold md:text-3xl">
            Sorry, we couldn't find this page.
          </p>
          <p className="mb-8 mt-4 text-gray-500 md:text-lg">
            But don't worry, you can find plenty of other things on our
            homepage.
          </p>
          <Link
            to="/"
            className="rounded bg-violet-600 px-8 py-3 font-semibold text-gray-50"
          >
            Back to Home
          </Link>
        </div>
      </div>
    </div>
  );
}
export default NotFound;
