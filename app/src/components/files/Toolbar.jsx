import React from "react";
import { NavLink } from "react-router-dom";
import { Toolbar } from "../Controls";

const FilesToolbar = ({
  path,
  file,
  fileURL,
  canSave,
  canEdit,
  onSave,
  onDelete,
  onUpload,
  onShare,
  onMkdir,
  onRmdir
}) => {
  const components = path.split("/").filter(item => item.length > 0);
  const breadcrumbs = (
    <>
      <NavLink to="/files" exact>
        .
      </NavLink>
      {components.map((item, idx) => (
        <NavLink
          key={idx}
          to={`/files/${components.slice(0, idx + 1).join("/")}`}
          exact
        >
          {`/${item}`}
        </NavLink>
      ))}
    </>
  );
  return (
    <Toolbar title={breadcrumbs}>
      {file && (
        <>
          {canEdit && canSave && (
            <a href="#save" onClick={onSave}>
              <i className="material-icons-round">save</i>
            </a>
          )}
          <a href={fileURL} download={file.name}>
            <i className="material-icons-round">download</i>
          </a>
          {canEdit && (
            <a href="#delete" onClick={onDelete}>
              <i className="material-icons-round">delete</i>
            </a>
          )}
        </>
      )}
      {!file && canEdit && (
        <>
          <a href="#upload" onClick={onUpload}>
            <i className="material-icons-round">upload</i>
          </a>
          <a href="#mkdir" onClick={onMkdir}>
            <i className="material-icons-round">create_new_folder</i>
          </a>
          {path !== "/" && (
            <a href="#rmdir" onClick={onRmdir}>
              <i className="material-icons-round">delete</i>
            </a>
          )}
          <a href="#share" onClick={onShare}>
            <i className="material-icons-round">share</i>
          </a>
        </>
      )}
    </Toolbar>
  );
};

export default FilesToolbar;
