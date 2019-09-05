import React from "react";
import styled from "@emotion/styled/macro";
import { connect } from "react-redux";
import { NavLink } from "react-router-dom";

const Badge = styled.small`
  flex: 0 0 auto;
  padding: 0.25em 0.5em;

  background: var(--outline-color);
  border-radius: 3px;

  font-weight: bold;
`;

const Albums = styled.ul`
  list-style: none;

  a {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  a span {
    margin-right: 1rem;

    overflow: hidden;
    white-space: nowrap;
    text-overflow: ellipsis;
  }

  ${Badge} {
    background: #e5e9f0;
  }

  a.active > ${Badge} {
    background: var(--outline-color);
  }
`;

export const AlbumsList = ({ albums, match }) => {
  const albumItems = albums.map(({ name, items_count }) => (
    <li key={name}>
      <NavLink to={`${match.url}/${name}`}>
        <span>{name}</span>

        <Badge>{items_count}</Badge>
      </NavLink>
    </li>
  ));

  return (
    <div>
      <h2>Albums</h2>
      <Albums>{albumItems}</Albums>
    </div>
  );
};

export default connect(
  ({ gallery: { albums } }) => ({ albums }),
  {}
)(AlbumsList);
