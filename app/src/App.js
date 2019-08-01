import React, { Component } from "react";
import styled from "@emotion/styled";
import { connect } from "react-redux";
import { HashRouter, Route, NavLink, Switch, Redirect } from "react-router-dom";

import AlbumsList from "./components/AlbumsList/index";
import LoginForm from "./components/LoginForm/index";
import { GalleryRoutes } from "./components/Routes";
import { apiClient, fetchModules } from "./actions";

const supportedModules = {
  gallery: {
    sidebar: AlbumsList,
    content: GalleryRoutes,
    icon: null
  }
};

const PageContainer = styled.div`
  min-height: 100vh;

  display: flex;
  background: #eceff4;
`;

const Sidepanel = styled.aside`
  padding: ${({ collapsed }) => (collapsed ? 0 : "5rem 0 5rem 2rem")};
  flex: 0 0 auto;
  max-width: ${({ collapsed }) => (collapsed ? 0 : 30)}%;

  nav {
    position: sticky;
    top: 2rem;
  }
`;

const Content = styled.main`
  flex: 1;
  margin: 2rem 2rem 2rem 2rem;
  padding: 4rem 3rem 1rem 3rem;

  background: white;
  border-radius: 8px;
  box-shadow: rgba(184, 194, 215, 0.25) 0px 4px 6px,
    rgba(184, 194, 215, 0.1) 0px 5px 7px;
`;

class App extends Component {
  componentDidMount() {
    apiClient.token = this.props.authToken;
    this.props.fetchModules();
  }

  componentDidUpdate(prevProps) {
    if (prevProps.authToken === this.props.authToken) return;
    apiClient.token = this.props.authToken;
    this.props.fetchModules();
  }

  _renderModules = () => {
    const { modules } = this.props;

    let navItems = [],
      sidebarItems = [],
      contentItems = [];
    modules.forEach(module => {
      const component = supportedModules[module];
      if (!component) return;

      navItems.push(
        <li key={module}>
          <NavLink to={module}>{component.icon}</NavLink>
        </li>
      );

      sidebarItems.push(
        <Route key={module} path={`/${module}`} component={component.sidebar} />
      );

      contentItems.push(
        <Route key={module} path={`/${module}`} component={component.content} />
      );
    });

    return { navItems, sidebarItems, contentItems };
  };

  render() {
    const { modules, authError } = this.props;
    const { navItems, sidebarItems, contentItems } = this._renderModules();

    return (
      <HashRouter>
        <PageContainer>
          <Sidepanel collapsed={navItems.length === 0}>
            {navItems.length > 1 && (
              <nav>
                <ul>{navItems}</ul>
              </nav>
            )}

            <nav>{sidebarItems}</nav>
          </Sidepanel>

          <Content>
            <Switch>
              <Route path="/login" component={LoginForm} />
              {authError && <Redirect to="/login" />}
              {contentItems}
              {modules.length > 0 && <Redirect to={`/${modules[0]}`} />}
              <Route render={() => <h2>No active modules</h2>} />
            </Switch>
          </Content>
        </PageContainer>
      </HashRouter>
    );
  }
}

export default connect(
  ({ modules, authError, authToken }) => ({ modules, authError, authToken }),
  { fetchModules }
)(App);
