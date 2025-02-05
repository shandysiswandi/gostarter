import { useDocumentTitle } from "@/hooks";

const Dashboard: React.FC = () => {
  useDocumentTitle("Dashboard");

  return (
    <>
      <h2 className="mb-6 text-3xl font-semibold">Welcome to my app</h2>
      <p>This is a basic layout using React, Tailwind CSS, and DaisyUI.</p>
    </>
  );
};

export default Dashboard;
