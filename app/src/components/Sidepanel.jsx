import React from "react";
import styled from "@emotion/styled/macro";
import { NavLink } from "react-router-dom";

const Container = styled.div`
  flex: 0 0 auto;
  position: sticky;
  height: 100vh;
  top: 0;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  max-width: ${({ collapsed }) => (collapsed ? 0 : 80)}%;
  padding: 4rem ${({ collapsed }) => (collapsed ? 1 : 2)}rem;
  margin-right: ${({ collapsed }) => (collapsed ? -2 : 0)}rem;

  nav {
    overflow: hidden;
  }

  @media (min-width: 700px) {
    max-width: ${({ collapsed }) => (collapsed ? 0 : 30)}%;
    margin-right: -0.5rem;
  }
`;

const SettingsBar = styled.div`
  flex: 0 0 auto;
  display: flex;
  justify-content: space-evenly;
  padding-top: 1rem;
  overflow: hidden;

  a {
    display: flex;
    align-items: center;
    margin-right: 1rem;

    color: var(--secondary-color);
  }

  a:last-of-type {
    margin-right: 0;
  }

  a small {
    margin-left: 0.5rem;
  }
`;

const Sidepanel = ({ children, collapsed, onSignOut }) => {
  return (
    <Container collapsed={collapsed}>
      <nav>{children}</nav>

      <SettingsBar>
        <NavLink to="/shares">
          <i className="material-icons-round">link</i>
          <small>Shares</small>
        </NavLink>
        <a href="#lougout" onClick={onSignOut}>
          <i className="material-icons-round">power_settings_new</i>
          <small>Logout</small>
        </a>
      </SettingsBar>
    </Container>
  );
};

export default Sidepanel;
