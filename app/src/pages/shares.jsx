import React, { useEffect } from "react";
import styled from "@emotion/styled";
import { connect } from "react-redux";
import { NavLink } from "react-router-dom";

import { apiClient, fetchShares, removeShare } from "../actions";

const Thumbs = styled.div`
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(150px, 1fr));
  grid-auto-rows: 100px;
  grid-gap: 1rem;
  margin-bottom: 1rem;

  div {
    overflow: hidden;
  }

  img {
    width: 100%;
    height: 100%;
    object-fit: cover;
  }
`;

export const GalleryItems = ({ gallery, items }) => {
  return (
    <Thumbs>
      {items.map((item, idx) => (
        <div key={idx}>
          <img src={apiClient.imageURL(gallery, item, "thumbnail")} alt="" />
        </div>
      ))}
    </Thumbs>
  );
};

const Share = styled.div`
  margin-bottom: 7rem;

  h3 {
    display: flex;
    align-items: center;

    a {
      display: flex;
    }
  }

  button {
    display: flex;
    align-items: center;
    padding: 0 1rem;
    margin-bottom: 0;
    background: var(--danger-color);
    border-color: var(--danger-color);
  }

  > div:last-of-type {
    display: flex;
    align-items: center;

    > div {
      display: flex;
      margin-left: 1rem;
    }

    i {
      margin-right: 0.5rem;
    }
  }
`;

const BackLink = styled(NavLink)`
  display: flex;
`;

export const SharesList = ({ shares, fetchShares, removeShare }) => {
  useEffect(() => {
    fetchShares();
  }, [fetchShares]);

  const shareItems = shares.map(({ slug, name, expires_at, items, type }) => (
    <Share key={slug}>
      <h3>
        <a href={`/share/${slug}`}>
          <i className="material-icons-round">link</i>
        </a>

        {name}
      </h3>

      {type === "gallery" && <GalleryItems gallery={name} items={items} />}

      <div>
        <button onClick={() => removeShare(slug)}>
          <i className="material-icons-round">delete</i> Remove
        </button>

        {expires_at && (
          <div>
            <i className="material-icons-round">access_time</i>
            {new Date(expires_at).toLocaleDateString()}
          </div>
        )}
      </div>
    </Share>
  ));

  return (
    <div>
      <h1>Shares</h1>
      {shareItems}
      {shares.length === 0 && <h2>No active shares</h2>}
      <BackLink to="/">
        <i className="material-icons-round">arrow_back_ios</i>
        Back
      </BackLink>
    </div>
  );
};

export default connect(
  ({ shares: { items } }) => ({ shares: items }),
  { fetchShares, removeShare }
)(SharesList);
