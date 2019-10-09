import React from "react";
import styled from "@emotion/styled/macro";
import { NavLink } from "react-router-dom";

const SettingsBar = styled.div`
  flex: 0 0 auto;
  display: flex;
  justify-content: space-evenly;
  padding-top: 1rem;

  a {
    display: flex;
    align-items: center;
    position: relative;
    margin-right: 1rem;

    color: var(--secondary-color);
  }

  a:last-of-type {
    margin-right: 0;
  }

  a:hover small {
    opacity: 1;
  }

  a small {
    position: absolute;
    top: -4rem;
    left: 0;
    padding: 0.5rem 1rem;
    transform: translate(-25%);
    background: white;
    border-radius: 5px;
    opacity: 0;
    transition: opasity 0.4s;
  }
`;

const Container = styled.div`
  flex: 0 0 auto;
  position: sticky;
  height: 100vh;
  top: 0;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  min-width: ${({ collapsed }) => (collapsed ? 0 : 250)}px;
  max-width: ${({ collapsed }) => (collapsed ? 0 : 80)}%;
  padding: 4rem ${({ collapsed }) => (collapsed ? 1 : 2)}rem;
  margin-right: ${({ collapsed }) => (collapsed ? -2 : 0)}rem;

  nav {
    width: ${({ collapsed }) => (collapsed ? 0 : "calc(100% + 5rem)")};
    overflow-x: hidden;
    overflow-y: auto;

    > div {
      margin-right: 5rem;
    }
  }

  ${SettingsBar} {
    overflow: ${({ collapsed }) => (collapsed ? "hidden" : "unset")};
  }

  @media (min-width: 700px) {
    max-width: ${({ collapsed }) => (collapsed ? 0 : 30)}%;
    margin-right: -0.5rem;
  }
`;

const Sidepanel = ({ children, collapsed, navItems, onSignOut }) => {
  return (
    <Container collapsed={collapsed}>
      <nav>{children}</nav>

      <SettingsBar>
        {navItems.map(({ title, route, icon }) => (
          <NavLink key={route} to={route}>
            <i className="material-icons-round">{icon}</i>
            <small>{title}</small>
          </NavLink>
        ))}
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
