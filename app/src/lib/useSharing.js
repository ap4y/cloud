import React, { useState, useEffect } from "react";
import styled from "@emotion/styled/macro";

import SharePopup from "../components/SharePopup";

const StickySharePopup = styled(SharePopup)`
  position: sticky;
  top: 1rem;
  width: 100%;
  margin: 0 auto 2rem auto;

  &:after {
    display: none;
  }

  @media (min-width: 800px) {
    max-width: 400px;
  }
`;

export default function useSharing(items, onShare) {
  const [sharedItems, setSharedItems] = useState([]);
  const [showSharing, setShowSharing] = useState(false);
  const [sharingSlug, setSharingSlug] = useState(null);
  const [sharingError, setSharingError] = useState(null);

  useEffect(() => {
    if (showSharing) {
      setSharedItems(items);
    } else {
      setSharingSlug(null);
      setSharingError(null);
    }
  }, [showSharing, items]);

  const toggleSharedItem = image => {
    if (!showSharing) return;

    if (sharedItems.includes(image)) {
      setSharedItems(sharedItems.filter(i => i !== image));
    } else {
      setSharedItems([...sharedItems, image]);
    }
  };

  const createShare = expireAt => {
    onShare(expireAt)
      .then(({ slug }) => setSharingSlug(slug))
      .catch(e => setSharingError(e.message));
  };

  const popup = (
    <StickySharePopup
      items={sharedItems}
      onShare={createShare}
      slug={sharingSlug}
      error={sharingError}
      onClose={() => setShowSharing(false)}
    />
  );

  return [popup, showSharing, setShowSharing, sharedItems, toggleSharedItem];
}
