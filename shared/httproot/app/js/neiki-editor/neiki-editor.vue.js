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
            //theme: 'dark',
            onChange: (content) => {
                self.$emit('update:modelValue', content)
            },
            onReady: () => {
                //Workaround: More-Menu is not configurable, hide Save & Autosave
                const moreMenu = self.$refs.wrapper.querySelector(".neiki-more-dropdown");
                for( const item of moreMenu.children) {
                    if(item.innerHTML.indexOf('Save') > -1 || item.innerHTML.indexOf('Autosave') > -1){
                        item.style = "display: none";
                    }
                }
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