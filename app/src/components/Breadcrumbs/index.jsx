import React, { Component } from 'react';
import { NavLink, withRouter } from 'react-router-dom';

import './styles.css';

class Breadcrumbs extends Component {
  render() {
    const pathComponents = this.props.location.pathname.split('/')
          .filter((item) => item.length > 0);
    const pathItems = pathComponents.map((item, idx) => {
      const items = pathComponents.slice(0, idx + 1);
      return (
        <li key={item}>
          <NavLink exact to={`/${items.join('/')}`}>{item}</NavLink>
        </li>
      );
    });
    return (
      <ul className="breadcrumbs">
        {pathItems}
      </ul>
    );
  }
}

export default withRouter(Breadcrumbs);
