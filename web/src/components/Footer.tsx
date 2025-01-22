const Footer: React.FC = () => {
  return (
    <footer className="bg-gray-200 p-4 text-center text-gray-700">
      © {new Date().getFullYear()} My Dashboard
    </footer>
  );
};

export default Footer;
