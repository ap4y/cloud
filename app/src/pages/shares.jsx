import React, { useEffect } from "react";
import styled from "@emotion/styled/macro";
import { connect } from "react-redux";

import { fetchShares, removeShare } from "../lib/actions";
import GalleryItems from "../components/shares/GalleryItems";
import FilesItems from "../components/shares/FilesItems";

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

const BackLink = styled.a`
  display: flex;
  padding: 1rem;
  background: white;
  color: var(--primary-text-color);
`;

const SharesContainer = styled.div`
  position: relative;

  ${BackLink} {
    position: absolute;
    top: -1rem;
    right: -1rem;
  }
`;

export const SharesList = ({ shares, fetchShares, removeShare, history }) => {
  useEffect(() => {
    fetchShares();
  }, [fetchShares]);

  const shareItems = shares.map(({ slug, name, expires_at, items, type }) => (
    <Share key={slug}>
      <h3>
        <a href={`/share/${slug}`}>
          <i className="material-icons-round">link</i>
        </a>

        {type === "gallery" && `Gallery: ${name}`}
        {type === "files" && `Folder: ${name}`}
      </h3>

      {type === "gallery" && <GalleryItems gallery={name} items={items} />}
      {type === "files" && <FilesItems folder={name} items={items} />}

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
    <SharesContainer>
      <h1>Shares</h1>
      <BackLink
        href="#back"
        onClick={e => {
          e.preventDefault();
          history.goBack();
        }}
      >
        <i className="material-icons-round">close</i>
      </BackLink>
      {shareItems}
      {shares.length === 0 && <h2>No active shares</h2>}
    </SharesContainer>
  );
};

export default connect(
  ({ shares: { items } }) => ({ shares: items }),
  { fetchShares, removeShare }
)(SharesList);
