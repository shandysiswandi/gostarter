const Footer: React.FC = () => {
  return (
    <footer className="bg-white p-4 text-center lg:ml-64 text-gray-700">
      © {new Date().getFullYear()} Go Starter
    </footer>
  );
};

export default Footer;
