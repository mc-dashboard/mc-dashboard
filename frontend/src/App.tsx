import { BrowserRouter, Routes, Route } from "react-router";
import Login from "./pages/Login";
import Dashboard from "./pages/Dashboard";
import { McTooltipPortal } from "./components/minecraft-ui";

function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<Login />} />
        <Route path="/dashboard" element={<Dashboard />} />
      </Routes>
      <McTooltipPortal />
    </BrowserRouter>
  );
}

export default App;
