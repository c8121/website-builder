import { selectAssetDialog } from '../asset-select-dialog.js'

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
            },
            onReady: () => {
                self.installImageDialog();
            }
        });

        // onReady may fire before self.editor is assigned, so try here as well
        self.installImageDialog();

        if (self.modelValue)
            self.editor.setContent(self.modelValue);
    },

    methods: {
        /**
         * Replaces the editor's built-in image dialog by the asset storage dialog.
         * Relies on internals of the editor: modal.open, restoreSavedSelection,
         * commands.insertImage, history.record and triggerChange. Re-check after updating the editor.
         */
        installImageDialog() {
            const self = this;

            const modal = self.editor && self.editor.modal;
            if (!modal || modal.assetDialogInstalled)
                return;
            modal.assetDialogInstalled = true;

            const openOriginal = modal.open.bind(modal);
            modal.open = (type, data) => {
                if (type !== 'image')
                    return openOriginal(type, data);

                // data.media is set if an existing image is edited
                const img = data && data.media;
                selectAssetDialog("Select image", (url) => {
                    if (img && img.tagName === 'IMG') {
                        img.src = url;
                        self.editor.history.record();
                        self.editor.triggerChange();
                    } else {
                        self.editor.restoreSavedSelection();
                        self.editor.commands.insertImage(url, '');
                    }
                });
            };
        },
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