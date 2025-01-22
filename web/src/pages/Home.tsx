import { Link } from "react-router";
import { useDocumentTitle } from "@/hooks";

const Home: React.FC = () => {
  useDocumentTitle("Home");

  return (
    <div className="flex h-screen items-center justify-center bg-gray-100">
      <Link
        to="/login"
        className="flex h-40 w-40 items-center justify-center rounded-full bg-violet-600 text-white hover:bg-violet-700 hover:shadow-md"
      >
        <h1 className="text-2xl font-bold">Let's Go</h1>
      </Link>
    </div>
  );
}

export default Home;