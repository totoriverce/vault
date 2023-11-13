# JsonEditor

`@container` is required if rendering within a modal - it gives context for the `<Hds::Copy::Button>` and sets `autoRefresh=true` so JsonEditor renders content (without this property @value only renders if editor is focused)

| Param            | Type                  | Description                                                                                                                                                                                                                                                                           |
| ---------------- | --------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [title]          | <code>string</code>   | Name above codemirror view                                                                                                                                                                                                                                                            |
| value            | <code>string</code>   | a specific string the comes from codemirror. It's the value inside the codemirror display                                                                                                                                                                                             |
| [valueUpdated]   | <code>function</code> | action to preform when you edit the codemirror value.                                                                                                                                                                                                                                 |
| [onFocusOut]     | <code>function</code> | action to preform when you focus out of codemirror.                                                                                                                                                                                                                                   |
| [helpText]       | <code>string</code>   | helper text.                                                                                                                                                                                                                                                                          |
| [extraKeys]      | <code>Object</code>   | Provides keyboard shortcut methods for things like saving on shift + enter.                                                                                                                                                                                                           |
| [gutters]        | <code>Array</code>    | An array of CSS class names or class name / CSS string pairs, each of which defines a width (and optionally a background), and which will be used to draw the background of the gutters.                                                                                              |
| [mode]           | <code>string</code>   | The mode defined for styling. Right now we only import ruby so mode must but be ruby or defaults to javascript. If you wanted another language you need to import it into the modifier.                                                                                               |
| [readOnly]       | <code>Boolean</code>  | Sets the view to readOnly, allowing for copying but no editing. It also hides the cursor. Defaults to false.                                                                                                                                                                          |
| [theme]          | <code>String</code>   | Specify or customize the look via a named "theme" class in scss.                                                                                                                                                                                                                      |
| [value]          | <code>String</code>   | Value within the display. Generally, a json string.                                                                                                                                                                                                                                   |
| [viewportMargin] | <code>String</code>   | Size of viewport. Often set to "Infinity" to load/show all text regardless of length.                                                                                                                                                                                                 |
| [example]        | <code>string</code>   | Example to show when value is null -- when example is provided a restore action will render in the toolbar to clear the current value and show the example after input                                                                                                                |
| [container]      | <code>string</code>   | Selector string or element object of containing element, set the focused element as the container value. This is for the Hds::Copy::Button and to set autoRefresh=true so content renders button [HDS docs](https://hds-website-hashicorp.vercel.app/components/copy/button?tab=code) |

**Example**

```hbs preview-template
<JsonEditor
  @title='Title here'
  @value={{stringify (hash foo='bar')}}
  @mode='ruby'
  @readOnly={{true}}
  @showToolbar={{true}}
/>
```
