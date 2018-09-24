import React, { Component } from 'react';
import { NavLink, Link } from 'react-router-dom';
import ReactSVG from 'react-svg';

import { apiClient } from '../../actions';

import './styles.css';

const meteringModes = {
  0: 'Unknown',
  1: 'Average',
  2: 'CenterWeightedAverage',
  3: 'Spot',
  4: 'MultiSpot',
  5: 'Pattern',
  6: 'Partial',
  255: 'other'
};

const flashModes = {
  0x0: 'Flash did not fire',
  0x1: 'Flash fired',
  0x5: 'Strobe return light not detected',
  0x7: 'Strobe return light detected',
  0x9: 'Flash fired, compulsory flash mode',
  0xD: 'Flash fired, compulsory flash mode, return light not detected',
  0xF: 'Flash fired, compulsory flash mode, return light detected',
  0x10: 'Flash did not fire, compulsory flash mode',
  0x18: 'Flash did not fire, auto mode',
  0x19: 'Flash fired, auto mode',
  0x1D: 'Flash fired, auto mode, return light not detected',
  0x1F: 'Flash fired, auto mode, return light detected',
  0x20: 'No flash function',
  0x41: 'Flash fired, red-eye reduction mode',
  0x45: 'Flash fired, red-eye reduction mode, return light not detected',
  0x47: 'Flash fired, red-eye reduction mode, return light detected',
  0x49: 'Flash fired, compulsory flash mode, red-eye reduction mode',
  0x4D: 'Flash fired, compulsory flash mode, red-eye reduction mode, return light not detected',
  0x4F: 'Flash fired, compulsory flash mode, red-eye reduction mode, return light detected',
  0x59: 'Flash fired, auto mode, red-eye reduction mode',
  0x5D: 'Flash fired, auto mode, return light not detected, red-eye reduction mode',
  0x5F: 'Flash fired, auto mode, return light detected, red-eye reduction mode'
};

const GalleryItem = (props) => {
  const { image, cellRef, path } = props;
  return (
    <li>
      <NavLink to={path}>
        <figure className="image-preview-cell" ref={cellRef}>
          <img alt={image.Name} src={`/api/gallery/thumbnails/${image.Path}`}/>
          <figcaption>{image.Name}</figcaption>
        </figure>
      </NavLink>
    </li>
  );
};

const EXIFData = (props) => {
  const { exif } = props;
  return (
    <ul className="image-preview-exif-data">
      <li><em>Created:</em> {exif.DateTime}</li>
      <li><em>Dimensions:</em> {exif.PixelXDimension}x{exif.PixelYDimension}</li>
      <li><em>Camera:</em> {exif.Make} {exif.Model}</li>
      <li><em>Lens:</em> {exif.LensModel}</li>
      <li><em>Exposure:</em> {exif.ExposureTime}</li>
      <li><em>Focal Length:</em> {exif.FNumber}</li>
      <li><em>ISO:</em> {exif.ISOSpeedRatings}</li>
      <li><em>Metering:</em> {meteringModes[exif.MeteringMode]}</li>
      <li><em>Flash:</em> {flashModes[exif.Flash]}</li>
    </ul>
  );
};

class ImagePreview extends Component {
  constructor(props) {
    super(props);
    this.cellRef = React.createRef();
    this.containerRef = React.createRef();
    this.previewRef = React.createRef();

    this.toggleFullscreen = this.toggleFullscreen.bind(this);
    this.toggleEXIF = this.toggleEXIF.bind(this);
    this.state = { fullscreen: false, exif: null };
  }

  componentDidUpdate() {
    this.centerCell();
  }

  centerCell() {
    const cell = this.cellRef.current;
    const container = this.containerRef.current;
    if (!cell || !container) return;

    container.scrollLeft = cell.offsetLeft - 0.5 * container.clientWidth + 0.5 * cell.clientWidth;
    cell.querySelector('img').onload = () => {
      container.scrollLeft = cell.offsetLeft - 0.5 * container.clientWidth + 0.5 * cell.clientWidth;
    };
  }

  imagePath(image) {
    return this.props.match.path.replace(':imageName', image.Name);
  }

  toggleFullscreen(e) {
    e.preventDefault();

    const fs = this.state.fullscreen,
          el = fs ? document : this.previewRef.current;

    var func;
    if (fs) {
      func = el.cancelFullScreen || el.webkitCancelFullScreen || el.mozCancelFullScreen || el.msCancelFullScreen;
    } else {
      func = el.requestFullScreen || el.webkitRequestFullScreen || el.mozRequestFullScreen || el.msRequestFullscreen;
    }


    if (func) {
      func.call(el);
      this.setState({ fullscreen: !fs });
    }
  }

  toggleEXIF(e) {
    e.preventDefault();

    if (this.state.exif) {
      this.setState({ exif: null });
    } else {
      const image = this.props.images.find((image) => image.Name === this.props.match.params.imageName);
      apiClient.do(`/api/gallery/exif/${image.Path}`).then(exif => { this.setState({ exif }); });
    }
  }

  render() {
    const { images, match } = this.props,
          selectedIdx = images.findIndex((image) => image.Name === match.params.imageName),
          selectedImage = images[selectedIdx],
          prevImage = images[selectedIdx - 1],
          nextImage = images[selectedIdx + 1];

    const galleryItems = this.props.images.map((image) => (
      <GalleryItem key={image.Name} image={image} path={this.imagePath(image)}
                   cellRef={selectedImage && image.Name === selectedImage.Name ? this.cellRef : React.createRef()}/>
    ));
    return (
      <div className="image-preview" ref={this.previewRef}>
        <header className="image-preview-header">
          {selectedImage &&
            <a href={`/api/gallery/images/${selectedImage.Path}`} download={selectedImage.Path.replace('/', '_')}>
              <ReactSVG className="icon" path="/images/ic_download.svg"/>
            </a>}
          <a href='#fullscreen' onClick={this.toggleFullscreen}>
            <ReactSVG className="icon" path="/images/ic_expand.svg"/>
          </a>
          <a href='#exif' onClick={this.toggleEXIF}>
            <ReactSVG className="icon" path="/images/ic_info_rect.svg"/>
          </a>
          <p>{this.props.galleryName}</p>
          <NavLink className="image-preview-close" exact to={match.path.replace('/:imageName', '')}>
            <ReactSVG className="icon" path="/images/ic_close.svg"/>
          </NavLink>
          {this.state.exif && <EXIFData exif={this.state.exif}/>}
        </header>

        {prevImage &&
          <Link to={this.imagePath(prevImage)} onClick={() => this.setState({ exif: null })} className="image-preview-nav-link prev">
            <ReactSVG className="icon" path="/images/ic_chevron_left.svg"/>
          </Link>}

        {selectedImage &&
          <img alt={selectedImage.Name} src={`/api/gallery/images/${selectedImage.Path}`}/>}        

        {nextImage &&
          <Link to={this.imagePath(nextImage)} onClick={() => this.setState({ exif: null })} className="image-preview-nav-link next">
            <ReactSVG className="icon" path="/images/ic_chevron_right.svg"/>
          </Link>}

        <ul className="image-preview-slides" ref={this.containerRef}>{galleryItems}</ul>
      </div>
    );
  }
}

export default ImagePreview;
