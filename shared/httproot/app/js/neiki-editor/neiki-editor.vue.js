export default {
    template: `
        <div ref="wrapper">
            <textarea ref="editor"></textarea>
        </div>
    `,

    props: {
        modelValue: {
            type: String,
            default: null
        },
    },

    mounted() {

        const self = this;
        self.editor = new NeikiEditor(self.$refs.editor, {
            themesUrl: '/js/neiki-editor/themes/',
            toolbar: [
                'viewCode', 'undo', 'redo', 'findReplace', '|',
                'bold', 'italic', 'underline', 'strikethrough', 'superscript', 'subscript', 'code', 'removeFormat', '|',
                'heading', /*'fontFamily',*/ 'fontSize', 'lineHeight', '|',
                'foreColor', 'backColor', '|',
                'alignLeft', 'alignCenter', 'alignRight', 'alignJustify', '|',
                'indent', 'outdent', '|',
                'bulletList', 'numberedList', 'blockquote', 'horizontalRule', '|',
                'insertDropdown'/*, '|',
                'themeToggle', 'moreMenu'*/
            ],
            onChange: (content) => {
                self.$emit('update:modelValue', content)
            }
        });

        if (self.modelValue)
            self.editor.setContent(self.modelValue);
    },

    unmounted() {
        const self = this;
        if (self.editor) {
            self.editor.destroy();
            console.log("Neikis editor unmounted");
        }
    },

    emits: ['update:modelValue'],
}