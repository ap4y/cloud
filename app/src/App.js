import React, { useState, useEffect } from "react";
import styled from "@emotion/styled/macro";
import { connect } from "react-redux";
import {
  BrowserRouter,
  Route,
  NavLink,
  Switch,
  Redirect
} from "react-router-dom";

import AlbumsList from "./pages/albums";
import LoginForm from "./pages/login";
import { ShareRoutes, GalleryRoutes } from "./Routes";
import SharesList from "./pages/shares";
import Sidepanel from "./components/Sidepanel";
import { apiClient, fetchModules, signOut } from "./actions";

const supportedModules = {
  gallery: {
    sidebar: AlbumsList,
    content: GalleryRoutes,
    icon: null
  }
};

const PageContainer = styled.div`
  min-height: 100vh;
  position: relative;

  display: flex;
  background: var(--primary-background-color);
`;

const Content = styled.main`
  flex: 1;
  position: relative;

  overflow-x: hidden;

  > div {
    padding: 2rem;

    background: white;
  }

  @media (min-width: 700px) {
    margin: 2rem;
    overflow-x: unset;

    > div {
      height: 100%;
      padding: 2rem 3rem;
      position: relative;
      z-index: 1;

      border-radius: 8px;
    }
  }
`;

const CollapseButton = styled.a`
  position: sticky;
  display: flex;
  align-items: center;
  padding: 1.5rem;
  margin-bottom: -3rem;
  top: 0;
  width: 100%;
  z-index: 1;
  cursor: pointer;

  background: white;

  color: var(--secondary-color);

  @media (min-width: 700px) {
    position: fixed;
    top: 5rem;
    margin-left: -3rem;
    margin-bottom: 0;
    padding: 0.5rem;
    width: auto;

    border-top-left-radius: 5px;
    border-bottom-left-radius: 5px;
    box-shadow: rgba(184, 194, 215, 0.25) 0px 4px 6px,
      rgba(184, 194, 215, 0.1) 0px 5px 7px;
  }
`;

const App = ({ modules, authError, authToken, fetchModules, signOut }) => {
  const [collapsed, setCollapsed] = useState(
    document.body.clientWidth >= 700 ? false : true
  );

  useEffect(() => {
    apiClient.token = authToken;
    fetchModules();
  }, [fetchModules, authToken]);

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

  const renderSidebar = () => (
    <Sidepanel collapsed={collapsed} onSignOut={signOut}>
      {sidebarItems}
    </Sidepanel>
  );

  const renderCollapseButton = () => (
    <CollapseButton onClick={() => setCollapsed(!collapsed)}>
      <i className="material-icons-round">{collapsed ? "menu" : "close"}</i>
    </CollapseButton>
  );

  return (
    <BrowserRouter>
      <PageContainer>
        <Switch>
          <Route path="/gallery" render={renderSidebar} />
        </Switch>
        <Content>
          <Route path="/gallery" render={renderCollapseButton} />

          <div>
            <Switch>
              <Route path="/share/:slug" component={ShareRoutes} />

              <Route path="/login" component={LoginForm} />
              {authError && <Redirect to="/login" />}

              {modules.length > 0 && (
                <Route path="/shares" component={SharesList} />
              )}

              {contentItems}
              {modules.length > 0 && <Redirect to={`/${modules[0]}`} />}

              <Route render={() => <h2>No active modules</h2>} />
            </Switch>
          </div>
        </Content>
      </PageContainer>
    </BrowserRouter>
  );
};

export default connect(
  ({ modules, authError, authToken }) => ({ modules, authError, authToken }),
  { fetchModules, signOut }
)(App);
