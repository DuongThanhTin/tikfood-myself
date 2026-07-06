import { render, screen } from "@testing-library/react";
import { describe, expect, it } from "vitest";

import { fallbackVenues } from "../lib/api";
import { VenueList } from "./VenueList";

describe("VenueList", () => {
  it("renders one card per venue with the result count", () => {
    render(<VenueList venues={fallbackVenues} />);

    expect(screen.getByText(`${fallbackVenues.length} venues`)).toBeInTheDocument();
    expect(screen.getByText("Banh Mi Hem")).toBeInTheDocument();
    expect(screen.getByText("Pho Bo Nguyen")).toBeInTheDocument();
  });

  it("shows the empty state when there are no venues", () => {
    render(<VenueList venues={[]} />);

    expect(screen.getByText("No venues match the current filters.")).toBeInTheDocument();
    expect(screen.getByText("0 venues")).toBeInTheDocument();
  });
});
