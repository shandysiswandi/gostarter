import { useDocumentTitle } from "@/hooks";

const Dashboard: React.FC = () => {
  useDocumentTitle("Dashboard");

  return (
    <>
      <h1 className="mb-4 text-2xl font-bold text-gray-800">Dashboard</h1>

      <div className="rounded-xl bg-white px-4 pb-4 pt-1 shadow">
        {[...Array(20).keys()].map((i) => (
          <div key={i} className="mt-4 rounded bg-gray-200 p-4">
            <h3 className="text-lg font-semibold">Section {i + 1}</h3>
            <p className="text-gray-700">
              This is the content of section {i + 1}.
            </p>
          </div>
        ))}
      </div>
    </>
  );
};

export default Dashboard;
