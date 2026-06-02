import { BrowserRouter, Routes, Route } from "react-router";
import Dashboard from "./pages/Dashboard";
import { McTooltipPortal } from "./components/minecraft-ui";

function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<Dashboard />} />
      </Routes>
      <McTooltipPortal />
    </BrowserRouter>
  );
}

export default App;
