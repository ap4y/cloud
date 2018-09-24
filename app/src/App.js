import React, { Component } from 'react';
import { HashRouter, Route, NavLink, Switch, Redirect } from 'react-router-dom';
import ReactSVG from 'react-svg';

import Breadcrumbs from './components/Breadcrumbs/index';
import AlbumsList from './components/AlbumsList/index';
import { GalleryRoutes } from './components/Routes';
import { apiClient } from './actions';

import './App.css';

class App extends Component {
  constructor(props) {
    super(props);
    this.state = {
      modules: [],
      moduleMappings: {
        gallery: { sidebar: AlbumsList, content: GalleryRoutes, icon: <ReactSVG path="/images/images-regular.svg"/> }
      }
    };
  }

  componentDidMount() {
    apiClient.do('/api/modules')
      .then(({ modules }) => { this.setState({ modules }); });
  }

  render() {
    let navItems = [], sidebarItems = [], contentItems = [], { modules } = this.state;
    modules.forEach((module) => {
      const component = this.state.moduleMappings[module];
      if (!component) return;

      navItems.push(
        <li key={module} className="app-module-nav-item">
          <NavLink to={module}>{component.icon}</NavLink>
        </li>
      );

      sidebarItems.push(
        <Route className="app-sidebar-items" key={module} path={`/${module}`} component={component.sidebar}/>
      );

      contentItems.push(
        <Route key={module} path={`/${module}`} component={component.content}/>
      );
    });

    let modulesNav = navItems.length <= 1 ? null : <nav className="app-modules"><ul>{navItems}</ul></nav>;

    return (
      <HashRouter>
        <div className="app">
          <header className="app-header">
            <Breadcrumbs/>
          </header>

          <aside className="app-sidebar">
            {modulesNav}

            <nav className="app-sidebar-nav">
              {sidebarItems}
            </nav>
          </aside>

          <main className="app-content">
            <Switch>
              {contentItems}
              {modules.length > 0 && <Redirect to={`/${modules[0]}`}/>}
              <Route render={() => (<h2>No active modules</h2>)}/>
            </Switch>
          </main>

          <footer className="app-footer"></footer>
        </div>
      </HashRouter>
    );
  }
}

export default App;
