import React, { useState, useRef, useEffect } from "react";
import styled from "@emotion/styled/macro";
import { NavLink, Link } from "react-router-dom";

import EXIFData from "./ImageEXIF";

import { apiClient } from "../actions";

const Figure = styled.figure`
  margin: 0;
  position: relative;
  display: flex;
  height: 100%;

  img {
    max-width: none;
    object-fit: cover;
  }

  figcaption {
    position: absolute;
    bottom: 0;
    right: 0;
    left: 0;
    text-align: center;
    background: var(--shade-color);
    opacity: ${({ selected }) => (selected ? 1 : 0)};
    transition: opacity 0.4s;
  }

  &:hover figcaption {
    opacity: 1;
  }
`;

export const AlbumItem = ({ image, src, selected }) => {
  return (
    <Figure selected={selected}>
      <img alt={image.name} src={src} />
      <figcaption>{image.name}</figcaption>
    </Figure>
  );
};

const Container = styled.div`
  position: fixed;
  top: 0;
  bottom: 0;
  right: 0;
  left: 0;
  z-index: 100;
  display: flex;
  flex-direction: column;

  background: var(--secondary-background-color);
`;

const Thumbs = styled.ul`
  flex: 0 0 auto;
  list-style: none;
  display: flex;
  margin: 0;
  padding: 0 0.5rem;
  max-height: ${({ hidden }) => (hidden ? 0 : 20)}%;
  overflow-x: auto;

  background: #3b4252;
  color: var(--primary-background-color);

  li {
    margin: 0;
    padding: 1rem 0.5rem;
  }

  a {
    display: block;
    height: 100%;
    color: inherit;
  }

  @media (min-width: 700px) {
    max-height: ${({ hidden }) => (hidden ? 0 : 15)}%;
  }
`;

const ImageContainer = styled.div`
  flex: 1;
  position: relative;

  img {
    position: absolute;
    top: 0;
    bottom: 0;
    left: 0;
    right: 0;
    width: 100%;
    height: 100%;
    object-fit: contain;
  }

  a {
    position: absolute;
    top: 50%;
    transform: translate(0, -50%);
    z-index: 101;
    opacity: 0;
    transition: opacity 0.4s;
    color: var(--primary-background-color);
  }

  a i {
    font-size: 96px;
  }

  a:first-of-type {
    left: 0%;
  }

  a:last-of-type {
    right: 0%;
  }

  &:hover a {
    opacity: 1;
  }
`;

const EXIFContainer = styled.div`
  position: absolute;
  left: 0;
  top: 5rem;
  z-index: 102;
  opacity: 0.7;
`;

const Toolbar = styled.header`
  display: flex;
  align-items: center;
  justify-content: space-between;
  height: 5rem;
  padding: 1rem;
  position: relative;
  z-index: 103;

  background: var(--secondary-background-color);
  box-shadow: rgba(36, 41, 51, 0.15) 0px 5px 10px 0px;
  color: var(--primary-background-color);

  & > div {
    display: flex;
  }

  h4 {
    margin: 0;
  }

  a {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 3rem;
    height: 3rem;
    margin-right: 1rem;

    color: var(--primary-background-color);
  }

  a:last-of-type {
    margin-right: 0;
  }

  @media (min-width: 700px) {
    h4 {
      position: absolute;
      top: 50%;
      left: 50%;
      transform: translate(-50%, -50%);
    }
  }
`;

const ImagePreview = ({ images, albumName, share, match, fetchExif }) => {
  const [fullscreen, setFullscreen] = useState(false);
  const [exif, setExif] = useState(null);
  const [selectedIdx, setSelectedIdx] = useState(0);

  const previewRef = useRef(null);
  const thumbsRef = useRef(null);

  useEffect(() => {
    document.body.style.overflow = "hidden";
    document.onfullscreenchange = e => {
      if (!document.fullscreenElement) setFullscreen(false);
    };

    return () => {
      document.body.style.overflow = null;
    };
  }, []);

  useEffect(() => {
    const selectedIdx = images.findIndex(
      image => image.name === match.params.imageName
    );
    setSelectedIdx(selectedIdx);
  }, [images, match.params.imageName]);

  useEffect(() => {
    const thumbs = thumbsRef.current;
    if (!thumbs) return;

    const cell = thumbs.children[selectedIdx];
    if (cell) cell.scrollIntoView({ inline: "center", behavior: "smooth" });
  }, [selectedIdx]);

  const selectedImage = images[selectedIdx];

  const imagePath = image => match.path.replace(":imageName", image.name);

  const toggleFullscreen = e => {
    e.preventDefault();

    const fs = fullscreen,
      el = fs ? document : previewRef.current;

    var func;
    if (fs) {
      func =
        el.cancelFullScreen ||
        el.webkitCancelFullScreen ||
        el.mozCancelFullScreen ||
        el.msCancelFullScreen;
    } else {
      func =
        el.requestFullScreen ||
        el.webkitRequestFullScreen ||
        el.mozRequestFullScreen ||
        el.msRequestFullscreen;
    }

    if (func) {
      func.call(el);
      setFullscreen(!fs);
    }
  };

  const toggleEXIF = e => {
    e.preventDefault();

    if (exif) {
      setExif(null);
      return;
    }

    fetchExif(albumName, selectedImage.path, share).then(exif => setExif(exif));
  };

  const imageURL = ({ path }, type = "image") =>
    apiClient.imageURL(albumName, path, type, share);

  const prevImage =
    images[selectedIdx === 0 ? images.length - 1 : selectedIdx - 1];
  const nextImage =
    images[selectedIdx === images.length - 1 ? 0 : selectedIdx + 1];

  const albumItems = images.map((image, idx) => (
    <li key={image.name}>
      <NavLink to={imagePath(image)}>
        <AlbumItem
          selected={idx === selectedIdx}
          image={image}
          src={imageURL(image, "thumbnail")}
        />
      </NavLink>
    </li>
  ));

  return (
    <Container ref={previewRef}>
      <Toolbar>
        <div>
          {selectedImage && (
            <a
              href={imageURL(selectedImage)}
              download={selectedImage.path.replace("/", "_")}
            >
              <i className="material-icons-round">get_app</i>
            </a>
          )}
          <a href="#exif" onClick={toggleEXIF}>
            <i className="material-icons-round">info</i>
          </a>
          <a href="#fullscreen" onClick={toggleFullscreen}>
            <i className="material-icons-round">
              {fullscreen ? "fullscreen_exit" : "fullscreen"}
            </i>
          </a>
        </div>

        <h4>{albumName}</h4>

        <NavLink exact to={match.path.replace("/:imageName", "")}>
          <i className="material-icons-round">close</i>
        </NavLink>
      </Toolbar>

      <EXIFContainer>{exif && <EXIFData exif={exif} />}</EXIFContainer>

      {selectedImage && (
        <ImageContainer>
          <Link to={imagePath(prevImage)} onClick={() => setExif(null)}>
            <i className="material-icons-round">chevron_left</i>
          </Link>

          <img alt={selectedImage.name} src={imageURL(selectedImage)} />

          <Link to={imagePath(nextImage)} onClick={() => setExif(null)}>
            <i className="material-icons-round">chevron_right</i>
          </Link>
        </ImageContainer>
      )}

      <Thumbs ref={thumbsRef} hidden={fullscreen}>
        {albumItems}
      </Thumbs>
    </Container>
  );
};

export default ImagePreview;
