import React, { Component } from "react";
import styled from "@emotion/styled";
import { connect } from "react-redux";
import { NavLink, Link } from "react-router-dom";
import ReactSVG from "react-svg";

import EXIFData from "./exif";

import { apiClient } from "../../actions";

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

const GalleryItem = ({ image, cellRef, authToken, selected }) => {
  return (
    <Figure ref={cellRef} selected={selected}>
      <img
        alt={image.Name}
        src={`/api/gallery/thumbnails/${image.Path}?jwt=${authToken}`}
      />
      <figcaption>{image.Name}</figcaption>
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

  svg {
    height: 150px;
    width: auto;
    fill: var(--primary-background-color);
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
    height: 3rem;
    margin-right: 1rem;
  }

  a:last-of-type {
    margin-right: 0;
  }

  svg {
    height: 3rem;
    width: auto;
    fill: var(--primary-background-color);
  }
`;

class ImagePreview extends Component {
  cellRef = React.createRef();
  previewRef = React.createRef();

  state = { fullscreen: false, exif: null };

  componentDidMount() {
    document.body.style.overflow = "hidden";
    document.onfullscreenchange = e => {
      if (!document.fullscreenElement) this.setState({ fullscreen: false });
    };
    window.setTimeout(() => this.centerCell(), 400);
  }

  componentWillUnmount() {
    document.body.style.overflow = null;
  }

  componentDidUpdate() {
    this.centerCell();
  }

  centerCell() {
    const cell = this.cellRef.current;
    if (!cell) return;

    cell.scrollIntoView({ inline: "center", behavior: "smooth" });
  }

  imagePath = image => this.props.match.path.replace(":imageName", image.Name);

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
    } else {
      const image = this.props.images.find(
        image => image.Name === this.props.match.params.imageName
      );
      apiClient.do(`/gallery/exif/${image.Path}`).then(exif => {
        this.setState({ exif });
      });
    }
  };

  render() {
    const { fullscreen } = this.state;
    const { images, match, authToken } = this.props;
    const selectedIdx = images.findIndex(
      image => image.Name === match.params.imageName
    );
    const selectedImage = images[selectedIdx];
    const prevImage =
      images[selectedIdx == 0 ? images.length - 1 : selectedIdx - 1];
    const nextImage =
      images[selectedIdx == images.length - 1 ? 0 : selectedIdx + 1];

    const galleryItems = images.map((image, idx) => (
      <li key={image.Name}>
        <NavLink to={this.imagePath(image)}>
          <GalleryItem
            selected={idx === selectedIdx}
            image={image}
            authToken={authToken}
            cellRef={
              selectedImage && image.Name === selectedImage.Name
                ? this.cellRef
                : React.createRef()
            }
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
                href={`/api/gallery/images/${selectedImage.Path}?jwt=${authToken}`}
                download={selectedImage.Path.replace("/", "_")}
              >
                <ReactSVG path="/images/ic_download.svg" />
              </a>
            )}
            <a href="#fullscreen" onClick={this.toggleFullscreen}>
              {!fullscreen && <ReactSVG path="/images/ic_fullscreen.svg" />}
              {fullscreen && <ReactSVG path="/images/ic_fullscreen_exit.svg" />}
            </a>
            <a href="#exif" onClick={this.toggleEXIF}>
              <ReactSVG path="/images/ic_info.svg" />
            </a>
          </div>

          <h4>{this.props.galleryName}</h4>

          <NavLink exact to={match.path.replace("/:imageName", "")}>
            <ReactSVG path="/images/ic_close.svg" />
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
              <ReactSVG path="/images/ic_chevron_left.svg" />
            </Link>

            <img
              alt={selectedImage.Name}
              src={`/api/gallery/images/${selectedImage.Path}?jwt=${authToken}`}
            />

            <Link
              to={this.imagePath(nextImage)}
              onClick={() => this.setState({ exif: null })}
            >
              <ReactSVG path="/images/ic_chevron_right.svg" />
            </Link>
          </ImageContainer>
        )}

        <Thumbs hidden={fullscreen}>{galleryItems}</Thumbs>
      </Container>
    );
  }
}

export default connect(
  ({ authToken }) => ({ authToken }),
  {}
)(ImagePreview);