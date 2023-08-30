# In 69 words
I often find myself trying to navigate level up with pressing `<backspace>`, it is close to `<enter>` and feels like a natural inverse for it. However there's a problem with just configuring `<backspace>` as a return universal keybinding: in this case whenever I need to input text in a popup it happens so that it is imposible to close popup: empty input field does not push keypress event further.

# Notes
[Quit actions module](https://github.com/jesseduffield/lazygit/blob/2b26b380b6ac40d3acfbf0446237990c11e5e0d2/pkg/gui/controllers/quit_actions.go#L51C26-L51C32) maybe a good place to start figuring out.

Also there should be difference when `quitOnTopLevelReturn` config is set to true:
- Poping-out with <esc> should quit from the top level,
- Poping-out with <backspace> should just get to the 1 window.

Possile way to configure this is:
```yaml
keybinding:
  universal:
    quit: 'q'
    return:
      - !TopLevel:always:quit
	<esc>
      - !TopLevel:never:quit
	!Input:whenEmpty:buble-up
	<backspace>
---
```
