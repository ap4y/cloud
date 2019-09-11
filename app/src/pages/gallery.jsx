import React, { useState, useEffect } from "react";
import styled from "@emotion/styled/macro";
import { connect } from "react-redux";
import { NavLink, Route } from "react-router-dom";
import VisibilitySensor from "react-visibility-sensor";

import ImagePreview from "../components/ImagePreview";
import SharePopup from "../components/SharePopup";
import { apiClient, fetchAlbum, shareAlbum, fetchExif } from "../actions";

const Figure = styled.figure`
  position: relative;
  display: flex;
  margin: 0;
  height: 200px;

  background: var(--outline-color);
  border-radius: 5px;
  overflow: hidden;

  img {
    flex: 1 0 auto;
    height: 200px;
    object-fit: cover;
  }

  figcaption {
    position: absolute;
    left: 0;
    right: 0;
    bottom: 0;

    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 0.25em 0.5em;

    background: var(--secondary-background-color);
    color: white;
  }

  figcaption span {
    flex: 1 0 auto;
    margin-right: 1rem;
  }

  figcaption small {
    overflow: hidden;
    white-space: nowrap;
    text-overflow: ellipsis;
  }
`;

export const ImageCell = ({ image: { name, path, updated_at }, src }) => {
  const [render, setRender] = useState(false);

  return (
    <VisibilitySensor
      partialVisibility
      onChange={visible => visible && setRender(true)}
    >
      <Figure>
        {render && <img alt={name} src={src} />}
        <figcaption>
          <span>{name}</span>
          <small>
            <time dateTime={updated_at}>
              {new Date(updated_at).toLocaleString()}
            </time>
          </small>
        </figcaption>
      </Figure>
    </VisibilitySensor>
  );
};

const StickySharePopup = styled(SharePopup)`
  position: sticky;
  top: 1rem;
  width: 100%;
  margin: 0 auto 2rem auto;

  &:after {
    display: none;
  }

  @media (min-width: 800px) {
    max-width: 50%;
  }
`;

const SortControl = styled.div`
  display: flex;

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

const Toolbar = styled.div`
  display: flex;
  flex-direction: column;
  align-items: baseline;
  justify-content: space-between;

  > div {
    position: relative;
    display: flex;
    margin-bottom: 2rem;

    > a {
      display: block;
      margin-left: 3rem;
      height: 2rem;
      color: var(--secondary-color);
    }
  }

  @media (min-width: 700px) {
    flex-direction: row;

    > div {
      margin-bottom: 0;
    }
  }
`;

const Images = styled.div`
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  grid-gap: 2rem;

  @media (min-width: 400px) {
    grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  }
`;

const ImageItem = styled.div`
  position: relative;

  &:before {
    display: ${({ active }) => (active ? "block" : "none")};
    content: "";
    position: absolute;
    top: 0;
    bottom: 0;
    left: 0;
    right: 0;
    z-index: 1;
    cursor: pointer;
  }

  &:after {
    display: ${({ active }) => (active ? "flex" : "none")};
    content: "${({ selected }) => (selected ? "✓" : "")}";
    position: absolute;
    top: 1rem;
    right: 1rem;
    width: 25px;
    height: 25px;
justify-content: center;
    align-items: center;
    z-index: 2;

    background: var(--outline-color);
    border-radius: 15px;

    color: var(--primary-text-color);
    font-weight: 700;
  }
`;

export const ImageGrid = ({
  images,
  share,
  albumName,
  match,
  fetchAlbum,
  fetchExif,
  shareAlbum
}) => {
  const [sharedItems, setSharedItems] = useState([]);
  const [showSharing, setShowSharing] = useState(false);
  const [sharingSlug, setSharingSlug] = useState(null);
  const [sharingError, setSharingError] = useState(null);
  const [sorting, setSorting] = useState({ field: "name", order: "up" });

  useEffect(() => {
    fetchAlbum(albumName, share);
  }, [albumName, share, fetchAlbum]);

  useEffect(() => {
    setShowSharing(false);
  }, [albumName]);

  const toggleSharing = e => {
    e.preventDefault();
    if (showSharing) {
      setShowSharing(false);
    } else {
      setSharedItems(images);
      setShowSharing(true);
    }
  };

  const createShare = expireAt => {
    shareAlbum(albumName, sharedItems, expireAt)
      .then(({ slug }) => setSharingSlug(slug))
      .catch(e => setSharingError(e.message));
  };

  const finalizeSharing = () => {
    setShowSharing(false);
    setSharingSlug(null);
    setSharingError(null);
  };

  const changeSorting = (e, field) => {
    e.preventDefault();

    if (field !== sorting.field) {
      setSorting({ ...sorting, field });
      return;
    }

    const order = sorting.order === "up" ? "down" : "up";
    setSorting({ ...sorting, order });
  };

  const toggleSharedItem = image => {
    if (!showSharing) return;

    if (sharedItems.includes(image)) {
      setSharedItems(sharedItems.filter(i => i !== image));
    } else {
      setSharedItems([...sharedItems, image]);
    }
  };

  const imageItems = images
    .sort((a, b) => {
      const up = sorting.order === "up" ? 1 : -1;
      const down = sorting.order === "up" ? -1 : 1;
      return a[sorting.field] > b[sorting.field] ? up : down;
    })
    .map(image => {
      return (
        <ImageItem
          key={image.name}
          active={showSharing}
          selected={sharedItems.includes(image)}
          onClick={() => toggleSharedItem(image)}
        >
          <NavLink to={`${match.url}/${image.name}`}>
            <ImageCell
              image={image}
              src={apiClient.imageURL(
                albumName,
                image.path,
                "thumbnail",
                share
              )}
            />
          </NavLink>
        </ImageItem>
      );
    });

  return (
    <div>
      <Toolbar>
        <h2>{albumName}</h2>
        <div>
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
          {!share && (
            <a href="#share" onClick={toggleSharing}>
              <i className="material-icons-round">share</i>
            </a>
          )}
        </div>
      </Toolbar>
      {showSharing && (
        <StickySharePopup
          items={sharedItems}
          onShare={createShare}
          slug={sharingSlug}
          error={sharingError}
          onClose={finalizeSharing}
        />
      )}
      <Images>{imageItems}</Images>
      <Route
        path={`${match.url}/:imageName`}
        render={props => (
          <ImagePreview
            images={images}
            albumName={albumName}
            share={share}
            fetchExif={fetchExif}
            {...props}
          />
        )}
      />
    </div>
  );
};

export default connect(
  ({ albumImages }, props) => {
    const { albumName, slug } = props.match.params;
    return {
      albumName,
      share: slug,
      images: albumImages[albumName] || []
    };
  },
  { fetchAlbum, shareAlbum, fetchExif }
)(ImageGrid);
