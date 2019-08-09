import React, { Component } from "react";
import styled from "@emotion/styled";
import { connect } from "react-redux";
import { NavLink } from "react-router-dom";

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

  a > small {
    background: #e5e9f0;
  }

  a.active > small {
    background: var(--outline-color);
  }
`;

const Badge = styled.small`
  flex: 0 0 auto;
  padding: 0.25em 0.5em;

  background: var(--outline-color);
  border-radius: 3px;

  font-weight: bold;
`;

class AlbumsList extends Component {
  render() {
    const albumItems = this.props.albums.map(album => (
      <li key={album.name}>
        <NavLink to={`${this.props.match.url}/${album.name}`}>
          <span>{album.name}</span>

          <Badge>{album.items_count}</Badge>
        </NavLink>
      </li>
    ));

    return (
      <div>
        <h2>Albums</h2>
        <Albums>{albumItems}</Albums>
      </div>
    );
  }
}

export default connect(
  ({ gallery: { albums } }) => ({ albums }),
  {}
)(AlbumsList);
