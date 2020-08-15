import React from "react";
import styled from "@emotion/styled/macro";
import { Toolbar } from "../Controls";

const SortControl = styled.div`
  display: flex;
  margin-right: 2rem;

  a {
    display: flex;
    align-items: center;
    margin-left: 4px;
    font-size: 1.1rem;
    font-weight: 700;
    letter-spacing: 0.1rem;
    text-transform: uppercase;

    i {
      margin-left: -4px;
      margin-right: -6px;
    }
  }

  a:first-of-type {
    margin-left: 0;
  }
`;

const GalleryToolbar = ({
  albumName,
  allowSharing,
  sorting,
  onSort,
  onShare
}) => {
  const changeSorting = (e, field) => {
    e.preventDefault();

    if (field !== sorting.field) {
      onSort({ ...sorting, field });
      return;
    }

    const order = sorting.order === "up" ? "down" : "up";
    onSort({ ...sorting, order });
  };

  return (
    <Toolbar title={albumName}>
      <SortControl>
        <i className="material-icons-round">sort</i>
        <a href="#name" onClick={e => changeSorting(e, "name")}>
          <i
            style={{
              visibility: sorting.field === "name" ? "" : "hidden"
            }}
            className="material-icons-round"
          >{`arrow_drop_${sorting.order}`}</i>
          Name
        </a>
        <a href="#date" onClick={e => changeSorting(e, "updated_at")}>
          <i
            style={{
              visibility: sorting.field === "updated_at" ? "" : "hidden"
            }}
            className="material-icons-round"
          >{`arrow_drop_${sorting.order}`}</i>
          Date
        </a>
      </SortControl>
      {allowSharing && (
        <a href="#share" onClick={onShare}>
          <i className="material-icons-round">share</i>
        </a>
      )}
    </Toolbar>
  );
};

export default GalleryToolbar;
