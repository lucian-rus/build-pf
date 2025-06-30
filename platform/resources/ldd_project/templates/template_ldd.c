#include <linux/init.h>
#include <linux/kernel.h>
#include <linux/module.h>
#include <linux/proc_fs.h>
// generic imports
#include "template_ldd.h"

// --------------------------------------------------
// prototypes
// --------------------------------------------------
/* interface function prototypes */
static int     device_open(struct inode *inode, struct file *file);
static int     device_release(struct inode *inode, struct file *file);
static ssize_t device_read(struct file *filp, char __user *buffer, size_t length, loff_t *offset);
static ssize_t device_write(struct file *filp, const char __user *buffer, size_t length, loff_t *offset);

// --------------------------------------------------
// globals
// --------------------------------------------------
static const struct proc_ops fops = {
    .proc_open    = device_open,
    .proc_release = device_release,
    .proc_read    = device_read,
    .proc_write   = device_write,
};

/* handler for proc entry */
static struct proc_dir_entry *proc_entry = NULL;

// --------------------------------------------------
// entry points
// --------------------------------------------------
static int __init template_gpio_init(void) {
    /* create entry in proc */
    proc_entry = proc_create(TEMPLATE_GPIO_DEVNAME, 0666, NULL, &fops);

#if defined(GENERIC_LDD_LOGGER_ENABLED)
    LOG_WARN("template-ldd: registered\n");
#endif
    return 0;
}

static void __exit template_gpio_exit(void) {
    unregister_chrdev(TEMPLATE_GPIO_MAJOR, TEMPLATE_GPIO_DEVNAME);

#if defined(GENERIC_LDD_LOGGER_ENABLED)
    LOG_WARN("template-ldd: unregistered\n");
#endif
}

module_init(template_gpio_init);
module_exit(template_gpio_exit);

// --------------------------------------------------
// implementation for interface functions
// --------------------------------------------------
static int device_open(struct inode *inode, struct file *file) {
#if defined(GENERIC_LDD_LOGGER_ENABLED)
    LOG_WARN("template-ldd: device open called\n");
#endif
    return 0;
}

static int device_release(struct inode *inode, struct file *file) {
#if defined(GENERIC_LDD_LOGGER_ENABLED)
    LOG_WARN("template-ldd: device close called\n");
#endif
    return 0;
}

static ssize_t device_read(struct file *filp, char __user *buffer, size_t length, loff_t *offset) {
#if defined(GENERIC_LDD_LOGGER_ENABLED)
    LOG_WARN("template-ldd: read called\n");
#endif

    return 0;
}

static ssize_t device_write(struct file *filp, const char __user *buffer, size_t length, loff_t *offset) {
#if defined(GENERIC_LDD_LOGGER_ENABLED)
    LOG_WARN("template-ldd: write called\n");
#endif

    return length;
}

// --------------------------------------------------
// implementation for helper/generic functions
// --------------------------------------------------
