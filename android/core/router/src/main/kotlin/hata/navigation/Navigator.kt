package hata.navigation


interface Navigator {

    fun goTo(route: Route)

    fun goBack()
}
