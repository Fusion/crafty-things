import * as React from "react"
import * as ReactDOM from 'react-dom'
import craftXIconSrc from "./craftx-icon.png"
import { CraftTextBlock, CraftTextRun } from "@craftdocs/craft-extension-api";
import * as CraftBlockInteractor from "./craftBlockInteractor";

const App: React.FC<{}> = () => {
  const isDarkMode = useCraftDarkMode();

  React.useEffect(() => {
    if (isDarkMode) {
      document.body.classList.add("dark");
    } else {
      document.body.classList.remove("dark");
    }
  }, [isDarkMode]);

  return <div style={{
    display: "flex",
    flexDirection: "column",
    alignItems: "center",
  }}>
    <img className="icon" src={craftXIconSrc} alt="CraftX logo" />
    <button className={`btn ${isDarkMode ? "dark" : ""}`} onClick={insertActionButton}>
      Link in Things
    </button>
  </div>;
}

function useCraftDarkMode() {
  const [isDarkMode, setIsDarkMode] = React.useState(false);

  React.useEffect(() => {
    craft.env.setListener(env => setIsDarkMode(env.colorScheme === "dark"));
  }, []);

  return isDarkMode;
}

function insertActionButton() {
  let openTasks = CraftBlockInteractor.getUncheckedTodoItemsFromCurrentPage();
  openTasks.then((blocks) => {
    return Promise.all(
      blocks
        .filter(
          (block): block is CraftTextBlock => block.type === "textBlock"
        )
        .map((block) => {

          const taskForm = CraftBlockInteractor.getAddTaskForm(block);
          if (taskForm.title.includes("⤴")) return // Already tagged
          /* Debugging
          const outBlock = craft.blockFactory.textBlock({
            content: taskForm.title
          });
          craft.dataApi.addBlocks([block]);
          */

          craft.httpProxy.fetch({
            url: "http://localhost:48484/form",
            method: "POST",
            headers: {"Content-Type": "application/json"},
            body: {type: "text", text: JSON.stringify(taskForm)}
          })
          .then((response) => {
            if (response.status == "success") {
              response.data.body?.text().then((id) => {
                let result: CraftTextRun[] = []
                result.push(block.content[0])
                result.push({text: " ["})
                result.push({
                  text: "⤴", link: { type: "url", url: "things:///show?id=" + id }
                })
                result.push({text: "]"})
                block.content = result
                craft.dataApi.updateBlocks([block])
              })
            }
          });
        })
    )
  });
}

export function initApp() {
  ReactDOM.render(<App />, document.getElementById('react-root'))
}
