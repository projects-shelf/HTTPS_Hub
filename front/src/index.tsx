import React from "react";
import { createRoot } from "react-dom/client";
import { BrowserRouter, Routes, Route } from "react-router-dom";

import "./index.css";

import HomePage from "./pages/home";

const container = document.querySelector("#root");
if (!container) {
	throw new Error("No root element found");
}
const root = createRoot(container);

root.render(
	<React.StrictMode>
		<BrowserRouter>
			<Routes>
				<Route path="/" element={<HomePage />} />
			</Routes>
		</BrowserRouter>
	</React.StrictMode>,
);
