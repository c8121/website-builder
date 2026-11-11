import AssetStorageMain from './asset-storage-main.vue.js'
import Dialog from './dialog.js'

/**
 * Opens the asset storage in a dialog.
 * Calls callback(url) with the URL of the selected asset and closes the dialog.
 */
export function selectAssetDialog(title, callback) {
    return fetch(new URL('./asset-storage-main.tpl.html', import.meta.url))
        .then(response => response.text())
        .then(html => {
            AssetStorageMain.template = html;

            const dialog = Dialog.create();
            dialog.setTitle(title);
            dialog.setContentComponent(AssetStorageMain, {
                downloadCallback: (url) => {
                    callback(url);
                    dialog.close();
                }
            });
        });
}
