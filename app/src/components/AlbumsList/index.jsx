import React, { Component } from 'react';
import { connect } from 'react-redux';
import { NavLink } from 'react-router-dom';
import ReactSVG from 'react-svg';

import './styles.css';

class AlbumsList extends Component {
  render() {
    const albumItems = this.props.albums
          .map((album) =>
               <li className="albums-list-item" key={album.Name}>
                 <ReactSVG className="icon" path="/images/album.svg"/>
                 <NavLink to={`${this.props.match.url}/${album.Name}`}>
                   {album.Name}
                 </NavLink>
                 <span>{album.ItemsCount}</span>
               </li>
              );
    return (
      <div>
        <h3>Albums</h3>
        <ul className="albums-list">{albumItems}</ul>
      </div>
    );
  }
}

function mapStateToProps(state) {
  return { albums: state.gallery.albums };
}

export default connect(mapStateToProps, {})(AlbumsList);
