import { Link, NavLink } from 'react-router-dom';

export default function Header() {
  return (
    <header className="border-b">
      <div className="mx-auto max-w-6xl px-4 py-4 flex items-center justify-between">
        <Link to="/" className="font-semibold text-lg">Pickleball</Link>

        <nav className="flex gap-4 text-sm">
          <NavLink to="/" className={({ isActive }) => (isActive ? 'font-semibold' : 'text-gray-600')}>
            Home
          </NavLink>
          <NavLink to="/login" className={({ isActive }) => (isActive ? 'font-semibold' : 'text-gray-600')}>
            Login
          </NavLink>
        </nav>
      </div>
    </header>
  );
}
