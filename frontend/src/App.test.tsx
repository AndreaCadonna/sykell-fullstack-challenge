import { render, screen } from "@testing-library/react";
import { describe, test, expect } from "vitest";
import App from "./App";

// Happy path tests for initial setup
describe("App Component", () => {
  test("renders web crawler dashboard title", () => {
    render(<App />);
    const titleElement = screen.getByText(/Web Crawler Dashboard/i);
    expect(titleElement).toBeInTheDocument();
  });

  test("shows initial scaffolding message", () => {
    render(<App />);
    const scaffoldingText = screen.getByText(/Initial scaffolding/i);
    expect(scaffoldingText).toBeInTheDocument();
  });

  test("displays setup complete message", () => {
    render(<App />);
    const setupText = screen.getByText(
      /React \+ TypeScript \+ Tailwind CSS setup complete/i
    );
    expect(setupText).toBeInTheDocument();
  });
});
