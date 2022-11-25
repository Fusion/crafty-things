# Crafty Things

A basic extension for Craft.do.

Install this extension, and as long as you have Things installed, clicking "Link in Things" will copy the unchecked tasks from the current page to Things 3. Links to both applications will be included.

I am quite likely to add to the meager feature set in the future.

Note: I am very grateful to [FlohGro](https://github.com/FlohGro-dev/Craftist) who figured out more of the Craft API than I will ever have to. I am even lifting some of his own README!

# Installing

1. Enable Craft eXtensions.
    - Mac: At the top left of the app, click your avatar, select Preferences, click 'Advanced'. Under Craft eXtensions, click the dropdown, and select 'Enabled'.
    - Web: At the top left of the web app, click your avatar, select Craft eXtensions, toggle 'Craft eXtensions' on.
2. Download the `.craftx` file from the latest [release](https://github.com/fusion/crafty-things/releases/)
3. Download the helper service from the same page. This helper is necessary because Craft extensions do not have a way of talking directly to Things 3.
4. Install the extension:
    - At the bottom of the right side bar, the eXtensions logo is now visible. Click the '+' sign, select the file you downloaded from the previous step, then click open.
5. Run the service binary (`craftythingshelper`) and it should automagically install itself as a lightweight service.
6. Done - you can now use the eXtension.
