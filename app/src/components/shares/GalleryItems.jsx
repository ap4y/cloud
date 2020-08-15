import React from "react";
import styled from "@emotion/styled/macro";

import { apiClient } from "../../lib/actions";

const Thumbs = styled.div`
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(150px, 1fr));
  grid-auto-rows: 100px;
  grid-gap: 1rem;
  margin-bottom: 1rem;

  div {
    overflow: hidden;
  }

  img {
    width: 100%;
    height: 100%;
    object-fit: cover;
  }
`;

const GalleryItems = ({ gallery, items }) => {
  return (
    <Thumbs>
      {items.map((item, idx) => (
        <div key={idx}>
          <img src={apiClient.imageURL(gallery, item, "thumbnail")} alt="" />
        </div>
      ))}
    </Thumbs>
  );
};

export default GalleryItems;
