import Footer from "@/components/Footer";
import Header from "@/components/Header";
import Sidebar from "@/components/Sidebar";
import { useDocumentTitle } from "@/hooks";

const Dashboard: React.FC = () => {
  useDocumentTitle("Dashboard");

  return (
    <div className="flex min-h-screen bg-gray-100">
      <Sidebar />

      <div className="flex flex-1 flex-col">
        <Header />

        <main className="p-4">
          <div className="rounded bg-white p-4 shadow">
            <h2 className="mb-4 text-2xl font-bold">
              Welcome to the Dashboard
            </h2>
            <p className="text-gray-700">
              Here is the main content of the dashboard.
            </p>
          </div>
        </main>

        <Footer />
      </div>
    </div>
  );
}

export default Dashboard;
