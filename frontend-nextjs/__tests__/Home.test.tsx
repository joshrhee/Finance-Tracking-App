// Testing library doc: https://testing-library.com/docs/react-testing-library/intro/

import { render, screen } from "@testing-library/react";
import Home from "@/app/page";

describe("Home", () => {
    it("should have Docs text", () => {
        render(<Home />); // ARRANGE

        const myElement = screen.getByText("Docs"); // ACTION

        expect(myElement).toBeInTheDocument(); // ASSERT
    });

    it('should contain the text "information"', () => {
        render(<Home />); // ARRANGE

        const myElement = screen.getByText(/information/i); // ACTION

        expect(myElement).toBeInTheDocument(); // ASSERT
    });

    it('should have a heading"', () => {
        render(<Home />); // ARRANGE

        const myElement = screen.getByRole("heading", {
            name: "Learn"
        }); // ACTION

        expect(myElement).toBeInTheDocument(); // ASSERT
    });
});
