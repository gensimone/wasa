package api

import (
	"net/http"
	"path/filepath"
)

func (rt *_router) Handler() http.Handler {
	// Serve media
	rt.router.ServeFiles(filepath.Join(rt.media, "*filepath"), http.Dir(rt.rootMedia))

	rt.router.GET("/context", rt.wrap(rt.getContextReply))
	rt.router.GET("/liveness", rt.liveness)

	// Groups.
	rt.router.POST("/groups", rt.authRequest(rt.createGroup))
	rt.router.POST("/groups/:groupId", rt.authRequest(rt.addToGroup))
	rt.router.GET("/groups/:groupId", rt.authRequest(rt.getGroup))
	rt.router.DELETE("/groups/:groupId", rt.authRequest(rt.deleteGroup))
	rt.router.PUT("/groups/:groupId/name", rt.authRequest(rt.setGroupName))
	rt.router.PUT("/groups/:groupId/photo", rt.authRequest(rt.setGroupPhoto))
	rt.router.DELETE("/groups/:groupId/photo", rt.authRequest(rt.deleteGroupPhoto))
	rt.router.GET("/groups/:groupId/members", rt.authRequest(rt.getMemberIds))
	rt.router.DELETE("/groups/:groupId/user", rt.authRequest(rt.leaveGroup))
	rt.router.DELETE("/groups/:groupId/user/:userId", rt.authRequest(rt.removeUser))

	// Conversations.
	rt.router.GET("/conversations", rt.authRequest(rt.getMyConversations))
	rt.router.GET("/conversations/:id", rt.authRequest(rt.getConversation))
	rt.router.POST("/conversations/:id/message", rt.authRequest(rt.sendMessage))
	rt.router.POST("/conversations/:id/fmessage", rt.authRequest(rt.forwardMessage))

	// Reactions.
	rt.router.POST("/reactions/:messageId", rt.authRequest(rt.addReaction))
	rt.router.DELETE("/reactions/:messageId", rt.authRequest(rt.deleteReaction))
	rt.router.GET("/reactions/:messageId", rt.authRequest(rt.getReactions))

	// Comments.
	rt.router.DELETE("/comments/:messageId", rt.authRequest(rt.uncommentMessage))
	rt.router.POST("/comments/:messageId", rt.authRequest(rt.commentMessage))

	// Messages.
	rt.router.PUT("/messages/:messageId/receipts", rt.authRequest(rt.setMessageStatusAsRead))
	rt.router.GET("/messages/:messageId/receipts", rt.authRequest(rt.getReceipts))
	rt.router.GET("/messages/:messageId", rt.authRequest(rt.getMessage))
	rt.router.DELETE("/messages/:messageId", rt.authRequest(rt.deleteMessage))

	// Users.
	rt.router.GET("/users", rt.authRequest(rt.getUsers))
	rt.router.GET("/users/:userId", rt.authRequest(rt.getUserById))
	rt.router.PUT("/users/:userId/name", rt.authRequest(rt.setMyUserName))
	rt.router.PUT("/users/:userId/photo", rt.authRequest(rt.setMyPhoto))
	rt.router.DELETE("/users/:userId/photo", rt.authRequest(rt.deleteMyPhoto))

	// Authentication.
	rt.router.POST("/session", rt.doLogin)

	return rt.router
}
