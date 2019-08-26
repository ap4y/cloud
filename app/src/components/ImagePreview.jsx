import React, { Component } from "react";
import styled from "@emotion/styled";
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

export const AlbumItem = ({ image, authToken, album, selected }) => {
  return (
    <Figure selected={selected}>
      <img
        alt={image.name}
        src={`/api/gallery/${album}/thumbnail/${image.path}?jwt=${authToken}`}
      />
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
  max-height: ${({ hidden }) => (hidden ? 0 : 15)}%;
  overflow-x: auto;

  background: #3b4252;
  color: var(--primary-background-color);

  li {
    margin: 0;
    padding: 1rem 0.5rem;
    display: block;
    height: 100%;
  }

  a {
    display: block;
    height: 100%;
    color: inherit;
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
    position: absolute;
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%);
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
`;

class ImagePreview extends Component {
  previewRef = React.createRef();
  thumbsRef = React.createRef();

  state = { fullscreen: false, exif: null, selectedIdx: 0 };

  componentDidMount() {
    document.body.style.overflow = "hidden";
    document.onfullscreenchange = e => {
      if (!document.fullscreenElement) this.setState({ fullscreen: false });
    };

    const { images, match } = this.props;
    const selectedIdx = images.findIndex(
      image => image.name === match.params.imageName
    );
    this.setState({ selectedIdx }, this.centerCell);
  }

  componentWillUnmount() {
    document.body.style.overflow = null;
  }

  componentDidUpdate(prevProps) {
    const { images, match } = this.props;

    if (
      images === prevProps.images &&
      match.params.imageName === prevProps.match.params.imageName
    )
      return;

    const selectedIdx = images.findIndex(
      image => image.name === match.params.imageName
    );
    this.setState({ selectedIdx }, this.centerCell);
  }

  centerCell() {
    const thumbs = this.thumbsRef.current;
    if (!thumbs) return;

    const { selectedIdx } = this.state;
    const cell = thumbs.children[selectedIdx];
    if (cell) cell.scrollIntoView({ inline: "center", behavior: "smooth" });
  }

  imagePath = image => this.props.match.path.replace(":imageName", image.name);

  toggleFullscreen = e => {
    e.preventDefault();

    const fs = this.state.fullscreen,
      el = fs ? document : this.previewRef.current;

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
      this.setState({ fullscreen: !fs });
    }
  };

  toggleEXIF = e => {
    e.preventDefault();

    if (this.state.exif) {
      this.setState({ exif: null });
      return;
    }

    const { images, albumName } = this.props;
    const image = images[this.state.selectedIdx];
    apiClient.do(`/gallery/${albumName}/exif/${image.path}`).then(exif => {
      this.setState({ exif });
    });
  };

  render() {
    const { fullscreen, selectedIdx } = this.state;
    const { images, match, authToken, albumName } = this.props;

    const selectedImage = images[selectedIdx];
    const prevImage =
      images[selectedIdx === 0 ? images.length - 1 : selectedIdx - 1];
    const nextImage =
      images[selectedIdx === images.length - 1 ? 0 : selectedIdx + 1];

    const albumItems = images.map((image, idx) => (
      <li key={image.name}>
        <NavLink to={this.imagePath(image)}>
          <AlbumItem
            selected={idx === selectedIdx}
            album={albumName}
            image={image}
            authToken={authToken}
          />
        </NavLink>
      </li>
    ));
    return (
      <Container ref={this.previewRef}>
        <Toolbar>
          <div>
            {selectedImage && (
              <a
                href={`/api/gallery/${albumName}/image/${selectedImage.path}?jwt=${authToken}`}
                download={selectedImage.path.replace("/", "_")}
              >
                <i className="material-icons-round">get_app</i>
              </a>
            )}
            <a href="#exif" onClick={this.toggleEXIF}>
              <i className="material-icons-round">info</i>
            </a>
            <a href="#fullscreen" onClick={this.toggleFullscreen}>
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

        <EXIFContainer>
          {this.state.exif && <EXIFData exif={this.state.exif} />}
        </EXIFContainer>

        {selectedImage && (
          <ImageContainer>
            <Link
              to={this.imagePath(prevImage)}
              onClick={() => this.setState({ exif: null })}
            >
              <i className="material-icons-round">chevron_left</i>
            </Link>

            <img
              alt={selectedImage.name}
              src={`/api/gallery/${albumName}/image/${selectedImage.path}?jwt=${authToken}`}
            />

            <Link
              to={this.imagePath(nextImage)}
              onClick={() => this.setState({ exif: null })}
            >
              <i className="material-icons-round">chevron_right</i>
            </Link>
          </ImageContainer>
        )}

        <Thumbs ref={this.thumbsRef} hidden={fullscreen}>
          {albumItems}
        </Thumbs>
      </Container>
    );
  }
}

export default ImagePreview;
