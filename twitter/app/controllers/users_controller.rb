class UsersController < ApplicationController
    def new
        @new_user = User.new(name: params[:name], handle: params[:handle])
        @new_credential = @new_user.build_credential(email: params[:email], password: params[:password])

        @result = {
            # params: params,
            "success" => false,
            "error" => {
            }
        }

        if @new_user.valid? && @new_credential.valid?
            @new_user.save
            @new_credential.save
            @result["success"] = true
            @result["user"] = @new_user

            session[:current_user_id] = @new_user.id
        else 
            puts @new_user.errors.messages
            puts @new_credential.errors.messages

            if @new_user.errors.messages.key?(:handle) && @new_credential.errors.messages.key?(:email)
                @result["error"]["message"] = "Handle and Email already exist!"
            elsif @new_user.errors.messages.key?(:handle)
                @result["error"]["message"] = "Handle already exists!"
            elsif @new_credential.errors.messages.key?(:email)
                @result["error"]["message"] = "Email already exists!"
            else
                @result["error"]["message"] = "Invalid Credentials!"
            end

            @result["success"] = false
        end

        render json: {result: @result}
    end

    def verify

        @credentials = Credential.find_by(email: params[:email], password: params[:password])

        @result = {
            # params: params,
            "success" => false,
            "error" => {
            }
        }

        if @credentials != nil
            @result["success"] = true
            @result["user"]  = @credentials.user
            
            session[:user] = nil
            session[:user] = @credentials.user
            # session[@credentials.user.id]
            @result["session_id"]  = session.id

        else
            @result["success"] = false
            @result["error"]["message"] = "Invalid credentials!"
        end

        @result[:session] = session
        render json: {result: @result}
    end

    def logout
        reset_session
        @result = {
            # params: params,
            "success" => true,
            "error" => {
            }
        }
        # render json: {result: @result}
        render json: {result: @result}
    end

    # def checkSession

    #     @result = {
    #         # params: params,
    #         "success" => true,
    #         "error" => {
    #         },
    #         "session_id" => session.id
    #     }
    #     puts session, params[:session_id]
    #     # render json: {result: @result}
    #     if session.id.to_s == params[:session_id].to_s
    #         @result["session_active"] = true
    #     else
    #         @result["session_active"] = false
    #     end
    #     render json: {result: @result}
    # end

    def userInfo
    end

end
