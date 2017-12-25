const NProgress = require('nprogress');
require('nprogress/nprogress.css');

/**
 * Initializes module.
 */
if (!document.getElementById('nprogress')) {
    NProgress.start();
}
// Hook to onDone event to bypass require new object creation.
// This allows nProgress stopping from another module and stop it's inner loop.
document.getElementById('nprogress').addEventListener('onDone', NProgress.done);
