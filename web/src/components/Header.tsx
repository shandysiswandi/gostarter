export default function Header() {
  return (
    <header className="bg-base-100 sticky top-0 z-50 py-2">
      <div className="container">
        <div className="navbar px-0">
          <div className="navbar-start">
            <div className="dropdown">
              <label
                tabIndex={0}
                className="btn btn-circle btn-primary mr-1 lg:hidden"
              >
                <i className="bi bi-list text-2xl"></i>
              </label>
              <ul
                tabIndex={0}
                className="menu-compact menu dropdown-content rounded-box bg-base-200 mt-1 w-52 p-2 shadow"
              >
                <li>
                  <a href="#!">Home</a>
                </li>
                <li>
                  <a href="#!">Services</a>
                </li>
                <li>
                  <a href="#!">About</a>
                </li>
                <li>
                  <a href="#!">Work</a>
                </li>
                <li>
                  <a href="#!">Case Study</a>
                </li>
              </ul>
            </div>
            <a className="btn btn-ghost text-2xl normal-case">daisyUI</a>
          </div>
          <div className="navbar-center hidden lg:flex">
            <ul className="menu menu-horizontal p-0 font-medium">
              <li>
                <a href="#!">Home</a>
              </li>
              <li>
                <a href="#!">Services</a>
              </li>
              <li>
                <a href="#!">About</a>
              </li>
              <li>
                <a href="#!">Work</a>
              </li>
              <li>
                <a href="#!">Case Study</a>
              </li>
            </ul>
          </div>
          <div className="navbar-end">
            <a href="#!">LOGIN</a>
          </div>
        </div>
      </div>
    </header>
  );
}
