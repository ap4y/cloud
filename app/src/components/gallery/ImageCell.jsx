import React, { useState } from "react";
import styled from "@emotion/styled/macro";
import VisibilitySensor from "react-visibility-sensor";

const Figure = styled.figure`
  position: relative;
  display: flex;
  margin: 0;
  height: 200px;

  background: var(--outline-color);
  border-radius: 5px;
  overflow: hidden;

  img {
    flex: 1 0 auto;
    height: 200px;
    object-fit: cover;
  }

  figcaption {
    position: absolute;
    left: 0;
    right: 0;
    bottom: 0;

    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 0.25em 0.5em;

    background: var(--secondary-background-color);
    color: white;
  }

  figcaption span {
    flex: 1 0 auto;
    margin-right: 1rem;
  }

  figcaption small {
    overflow: hidden;
    white-space: nowrap;
    text-overflow: ellipsis;
  }
`;

const ImageCell = ({ image: { name, updated_at }, src }) => {
  const [render, setRender] = useState(false);

  return (
    <VisibilitySensor
      partialVisibility
      onChange={visible => visible && setRender(true)}
    >
      <Figure>
        {render && <img alt={name} src={src} />}
        <figcaption>
          <span>{name}</span>
          <small>
            <time dateTime={updated_at}>
              {new Date(updated_at).toLocaleString()}
            </time>
          </small>
        </figcaption>
      </Figure>
    </VisibilitySensor>
  );
};

export default ImageCell;
